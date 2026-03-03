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

type encrypt struct {
	*flags.VirtualMachineFlag
	kmsProvider string
	kmsKey      string
}

func init() {
	cli.Register("vm.encrypt", &encrypt{})
}

func (cmd *encrypt) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.VirtualMachineFlag, ctx = flags.NewVirtualMachineFlag(ctx)
	cmd.VirtualMachineFlag.Register(ctx, f)
	f.StringVar(&cmd.kmsProvider, "kms-provider", "", "KMS provider ID")
	f.StringVar(&cmd.kmsKey, "key-id", "", "KMS key ID")
}

func (cmd *encrypt) Description() string {
	return `Encrypt VM.

The VM must be powered off.
If -kms-provider is omitted the cluster default KMS provider is used.
If -key-id is omitted vSphere auto-generates a new key.

Examples:
  govc vm.encrypt -vm $vm
  govc vm.encrypt -vm $vm -kms-provider my-kp
  govc vm.encrypt -vm $vm -kms-provider my-kp -key-id <uuid>`
}

func (cmd *encrypt) Run(ctx context.Context, f *flag.FlagSet) error {
	vm, err := cmd.VirtualMachine()
	if err != nil {
		return err
	}
	if vm == nil {
		return flag.ErrHelp
	}

	cryptoKeyId := types.CryptoKeyId{
		KeyId: cmd.kmsKey,
	}
	if cmd.kmsProvider != "" {
		cryptoKeyId.ProviderId = &types.KeyProviderId{Id: cmd.kmsProvider}
	}

	spec := types.VirtualMachineConfigSpec{
		Crypto: &types.CryptoSpecEncrypt{
			CryptoKeyId: cryptoKeyId,
		},
	}

	task, err := vm.Reconfigure(ctx, spec)
	if err != nil {
		return err
	}
	return task.Wait(ctx)
}
