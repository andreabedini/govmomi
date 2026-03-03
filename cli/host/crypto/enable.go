// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package crypto

import (
	"context"
	"flag"

	"github.com/vmware/govmomi/cli"
	"github.com/vmware/govmomi/cli/flags"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type enable struct {
	*flags.HostSystemFlag
	provider  string
	keyID     string
	algorithm string
	keyData   string
}

func init() {
	cli.Register("host.crypto.enable", &enable{})
}

func (cmd *enable) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.HostSystemFlag, ctx = flags.NewHostSystemFlag(ctx)
	cmd.HostSystemFlag.Register(ctx, f)
	f.StringVar(&cmd.provider, "kms-provider", "", "KMS provider ID")
	f.StringVar(&cmd.keyID, "key-id", "", "Key ID for the initial host encryption key")
	f.StringVar(&cmd.algorithm, "algorithm", "", "Key algorithm (e.g. AES-256)")
	f.StringVar(&cmd.keyData, "key-data", "", "Hex-encoded key material")
}

func (cmd *enable) Process(ctx context.Context) error {
	return cmd.HostSystemFlag.Process(ctx)
}

func (cmd *enable) Description() string {
	return `Fully enable host crypto with an initial key.

The initial key is used for core-dump encryption.  -key-id, -algorithm and
-key-data are required; -kms-provider is optional.

Examples:
  govc host.crypto.enable -host esxi01 -kms-provider my-kp -key-id <uuid> -algorithm AES-256 -key-data <hex>`
}

func (cmd *enable) Run(ctx context.Context, f *flag.FlagSet) error {
	if cmd.keyID == "" || cmd.algorithm == "" || cmd.keyData == "" {
		return flag.ErrHelp
	}

	host, err := cmd.HostSystem()
	if err != nil {
		return err
	}

	var hs mo.HostSystem
	err = host.Properties(ctx, host.Reference(), []string{"configManager.cryptoManager"}, &hs)
	if err != nil {
		return err
	}
	if hs.ConfigManager.CryptoManager == nil {
		return object.ErrNotSupported
	}

	keyID := types.CryptoKeyId{KeyId: cmd.keyID}
	if cmd.provider != "" {
		keyID.ProviderId = &types.KeyProviderId{Id: cmd.provider}
	}

	keyData, err := hexToKeyData(cmd.keyData)
	if err != nil {
		return err
	}

	req := types.CryptoManagerHostEnable{
		This: *hs.ConfigManager.CryptoManager,
		InitialKey: types.CryptoKeyPlain{
			KeyId:     keyID,
			Algorithm: cmd.algorithm,
			KeyData:   keyData,
		},
	}

	_, err = methods.CryptoManagerHostEnable(ctx, host.Client(), &req)
	return err
}
