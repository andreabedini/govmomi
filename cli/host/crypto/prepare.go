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

type prepare struct {
	*flags.HostSystemFlag
}

func init() {
	cli.Register("host.crypto.prepare", &prepare{})
}

func (cmd *prepare) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.HostSystemFlag, ctx = flags.NewHostSystemFlag(ctx)
	cmd.HostSystemFlag.Register(ctx, f)
}

func (cmd *prepare) Process(ctx context.Context) error {
	return cmd.HostSystemFlag.Process(ctx)
}

func (cmd *prepare) Description() string {
	return `Put host into the crypto "prepared" state.

The prepared state makes the host ready to receive encryption keys from
vCenter without fully enabling host crypto.

Examples:
  govc host.crypto.prepare
  govc host.crypto.prepare -host esxi01`
}

func (cmd *prepare) Run(ctx context.Context, f *flag.FlagSet) error {
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

	req := types.CryptoManagerHostPrepare{
		This: *hs.ConfigManager.CryptoManager,
	}
	_, err = methods.CryptoManagerHostPrepare(ctx, host.Client(), &req)
	return err
}
