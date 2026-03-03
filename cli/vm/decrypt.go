// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package vm

import (
	"context"
	"flag"

	"github.com/vmware/govmomi/cli"
	"github.com/vmware/govmomi/cli/flags"
	"github.com/vmware/govmomi/vim25/types"
)

type decrypt struct {
	*flags.VirtualMachineFlag
}

func init() {
	cli.Register("vm.decrypt", &decrypt{})
}

func (cmd *decrypt) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.VirtualMachineFlag, ctx = flags.NewVirtualMachineFlag(ctx)
	cmd.VirtualMachineFlag.Register(ctx, f)
}

func (cmd *decrypt) Description() string {
	return `Decrypt VM.

The VM must be powered off.

Examples:
  govc vm.decrypt -vm $vm`
}

func (cmd *decrypt) Run(ctx context.Context, f *flag.FlagSet) error {
	vm, err := cmd.VirtualMachine()
	if err != nil {
		return err
	}
	if vm == nil {
		return flag.ErrHelp
	}

	spec := types.VirtualMachineConfigSpec{
		Crypto: &types.CryptoSpecDecrypt{},
	}

	task, err := vm.Reconfigure(ctx, spec)
	if err != nil {
		return err
	}
	return task.Wait(ctx)
}
