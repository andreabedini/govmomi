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

type disable struct {
	*flags.HostSystemFlag
}

func init() {
	cli.Register("host.crypto.disable", &disable{})
}

func (cmd *disable) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.HostSystemFlag, ctx = flags.NewHostSystemFlag(ctx)
	cmd.HostSystemFlag.Register(ctx, f)
}

func (cmd *disable) Process(ctx context.Context) error {
	return cmd.HostSystemFlag.Process(ctx)
}

func (cmd *disable) Description() string {
	return `Disable host crypto.

Examples:
  govc host.crypto.disable
  govc host.crypto.disable -host esxi01`
}

func (cmd *disable) Run(ctx context.Context, f *flag.FlagSet) error {
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

	req := types.CryptoManagerHostDisable{
		This: *hs.ConfigManager.CryptoManager,
	}
	_, err = methods.CryptoManagerHostDisable(ctx, host.Client(), &req)
	return err
}
