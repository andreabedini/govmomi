// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package key

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/vmware/govmomi/cli"
	"github.com/vmware/govmomi/cli/flags"
	"github.com/vmware/govmomi/crypto"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
)

type add struct {
	*flags.ClientFlag
	provider  string
	keyID     string
	algorithm string
	keyData   string
}

func init() {
	cli.Register("kms.key.add", &add{})
}

func (cmd *add) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.ClientFlag, ctx = flags.NewClientFlag(ctx)
	cmd.ClientFlag.Register(ctx, f)
	f.StringVar(&cmd.provider, "kms-provider", "", "KMS provider ID")
	f.StringVar(&cmd.keyID, "key-id", "", "Key ID to register")
	f.StringVar(&cmd.algorithm, "algorithm", "", "Key algorithm (e.g. AES-256)")
	f.StringVar(&cmd.keyData, "key-data", "", "Hex-encoded key material")
}

func (cmd *add) Description() string {
	return `Import a pre-existing crypto key into the vCenter CryptoManager.

Unlike kms.key.create (which asks the KMS to generate a new key), this
command imports key material that already exists (e.g. from a backup or an
external KMS).

Examples:
  govc kms.key.add -kms-provider my-kp -key-id abc123 -algorithm AES-256 -key-data <hex>`
}

func (cmd *add) Run(ctx context.Context, f *flag.FlagSet) error {
	if cmd.keyID == "" || cmd.algorithm == "" || cmd.keyData == "" {
		return flag.ErrHelp
	}

	c, err := cmd.Client()
	if err != nil {
		return err
	}

	m, err := crypto.GetManagerKmip(c)
	if err != nil {
		return err
	}

	keyID := types.CryptoKeyId{KeyId: cmd.keyID}
	if cmd.provider != "" {
		keyID.ProviderId = &types.KeyProviderId{Id: cmd.provider}
	}

	raw, err := hex.DecodeString(cmd.keyData)
	if err != nil {
		return fmt.Errorf("-key-data: %w (expected hex-encoded key material)", err)
	}
	keyData := base64.StdEncoding.EncodeToString(raw)

	req := types.AddKey{
		This: m.Reference(),
		Key: types.CryptoKeyPlain{
			KeyId:     keyID,
			Algorithm: cmd.algorithm,
			KeyData:   keyData,
		},
	}

	_, err = methods.AddKey(ctx, c, &req)
	return err
}
