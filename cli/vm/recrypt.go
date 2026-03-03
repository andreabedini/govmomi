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

type recrypt struct {
	*flags.VirtualMachineFlag
	kmsProvider string
	kmsKey      string
	deep        bool
}

func init() {
	cli.Register("vm.recrypt", &recrypt{})
}

func (cmd *recrypt) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.VirtualMachineFlag, ctx = flags.NewVirtualMachineFlag(ctx)
	cmd.VirtualMachineFlag.Register(ctx, f)
	f.StringVar(&cmd.kmsProvider, "kms-provider", "", "KMS provider ID (default: cluster default)")
	f.StringVar(&cmd.kmsKey, "key-id", "", "KMS key ID (default: provider generates a new key)")
	f.BoolVar(&cmd.deep, "deep", false, "Deep re-encrypt (re-encrypts all data; VM must be powered off with no snapshots)")
}

func (cmd *recrypt) Description() string {
	return `Re-encrypt VM with a (new) key.

Shallow re-encrypt (default) rotates only the metadata key and does not
require the VM to be powered off.  Deep re-encrypt rotates all data and
requires the VM to be powered off with no snapshots.

Examples:
  govc vm.recrypt -vm $vm -kms-provider my-kp
  govc vm.recrypt -vm $vm -kms-provider my-kp -deep`
}

func (cmd *recrypt) Run(ctx context.Context, f *flag.FlagSet) error {
	vm, err := cmd.VirtualMachine()
	if err != nil {
		return err
	}
	if vm == nil {
		return flag.ErrHelp
	}

	keyID := types.CryptoKeyId{
		KeyId: cmd.kmsKey,
	}
	if cmd.kmsProvider != "" {
		keyID.ProviderId = &types.KeyProviderId{Id: cmd.kmsProvider}
	}

	var cryptoSpec types.BaseCryptoSpec
	if cmd.deep {
		cryptoSpec = &types.CryptoSpecDeepRecrypt{NewKeyId: keyID}
	} else {
		cryptoSpec = &types.CryptoSpecShallowRecrypt{NewKeyId: keyID}
	}

	spec := types.VirtualMachineConfigSpec{
		Crypto: cryptoSpec,
	}

	task, err := vm.Reconfigure(ctx, spec)
	if err != nil {
		return err
	}
	return task.Wait(ctx)
}
