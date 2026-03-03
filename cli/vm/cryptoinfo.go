// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package vm

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/vmware/govmomi/cli"
	"github.com/vmware/govmomi/cli/flags"
	"github.com/vmware/govmomi/vim25/mo"
)

type cryptoinfo struct {
	*flags.VirtualMachineFlag
	*flags.OutputFlag
}

func init() {
	cli.Register("vm.crypto.info", &cryptoinfo{})
}

func (cmd *cryptoinfo) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.VirtualMachineFlag, ctx = flags.NewVirtualMachineFlag(ctx)
	cmd.VirtualMachineFlag.Register(ctx, f)
	cmd.OutputFlag, ctx = flags.NewOutputFlag(ctx)
	cmd.OutputFlag.Register(ctx, f)
}

func (cmd *cryptoinfo) Process(ctx context.Context) error {
	if err := cmd.VirtualMachineFlag.Process(ctx); err != nil {
		return err
	}
	return cmd.OutputFlag.Process(ctx)
}

func (cmd *cryptoinfo) Description() string {
	return `Display VM crypto state.

CryptoState values: unlocked, locked.

Examples:
  govc vm.crypto.info -vm $vm
  govc vm.crypto.info -vm $vm -json`
}

type vmCryptoInfoResult struct {
	CryptoState string `json:"cryptoState"`
}

func (r *vmCryptoInfoResult) Write(w io.Writer) error {
	fmt.Fprintln(w, r.CryptoState)
	return nil
}

func (cmd *cryptoinfo) Run(ctx context.Context, f *flag.FlagSet) error {
	vm, err := cmd.VirtualMachine()
	if err != nil {
		return err
	}
	if vm == nil {
		return flag.ErrHelp
	}

	var obj mo.VirtualMachine
	err = vm.Properties(ctx, vm.Reference(), []string{"runtime.cryptoState"}, &obj)
	if err != nil {
		return err
	}

	return cmd.WriteResult(&vmCryptoInfoResult{CryptoState: obj.Runtime.CryptoState})
}
