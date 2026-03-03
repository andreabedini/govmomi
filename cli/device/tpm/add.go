// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package tpm

import (
	"context"
	"flag"
	"fmt"

	"github.com/vmware/govmomi/cli"
	"github.com/vmware/govmomi/cli/flags"
	"github.com/vmware/govmomi/crypto"
	"github.com/vmware/govmomi/vim25/types"
)

type add struct {
	*flags.VirtualMachineFlag
	kmsProvider string
}

func init() {
	cli.Register("device.tpm.add", &add{})
}

func (cmd *add) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.VirtualMachineFlag, ctx = flags.NewVirtualMachineFlag(ctx)
	cmd.VirtualMachineFlag.Register(ctx, f)
	f.StringVar(&cmd.kmsProvider, "kms-provider", "", "KMS provider ID (default: cluster default)")
}

func (cmd *add) Description() string {
	return `Add Trusted Platform Module (TPM) device to VM.

The VM must be encrypted; if it is not already, it will be encrypted using
the specified KMS provider (or the cluster default if -kms-provider is omitted).

Examples:
  govc device.tpm.add -vm $vm
  govc device.tpm.add -vm $vm -kms-provider my-kp
  govc device.info tpm-*`
}

func (cmd *add) Run(ctx context.Context, f *flag.FlagSet) error {
	vm, err := cmd.VirtualMachine()
	if err != nil {
		return err
	}

	if vm == nil {
		return flag.ErrHelp
	}

	providerID := cmd.kmsProvider
	if providerID == "" {
		m, err := crypto.GetManagerKmip(vm.Client())
		if err != nil {
			return err
		}
		providerID, err = m.GetDefaultKmsClusterID(ctx, nil, true)
		if err != nil {
			return err
		}
		if providerID == "" {
			return fmt.Errorf("no KMS provider configured; use -kms-provider or configure a default KMS cluster")
		}
	}

	spec := types.VirtualMachineConfigSpec{
		Crypto: &types.CryptoSpecEncrypt{
			CryptoKeyId: types.CryptoKeyId{
				ProviderId: &types.KeyProviderId{Id: providerID},
			},
		},
		DeviceChange: []types.BaseVirtualDeviceConfigSpec{
			&types.VirtualDeviceConfigSpec{
				Device:    &types.VirtualTPM{},
				Operation: types.VirtualDeviceConfigSpecOperationAdd,
			},
		},
	}

	task, err := vm.Reconfigure(ctx, spec)
	if err != nil {
		return err
	}
	if err = task.Wait(ctx); err != nil {
		return err
	}

	// output name of device we just created
	devices, err := vm.Device(ctx)
	if err != nil {
		return err
	}

	d := &types.VirtualTPM{}
	devices = devices.SelectByType(d)

	name := devices.Name(devices[len(devices)-1])

	fmt.Println(name)

	return nil
}
