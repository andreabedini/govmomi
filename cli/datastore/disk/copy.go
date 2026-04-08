// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package disk

import (
	"context"
	"flag"
	"fmt"

	"github.com/vmware/govmomi/cli"
	"github.com/vmware/govmomi/cli/flags"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

type cp struct {
	*flags.DatastoreFlag

	adapterType string
	diskType    string
	force       bool
}

func init() {
	cli.Register("datastore.disk.copy", &cp{})
}

func (cmd *cp) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.DatastoreFlag, ctx = flags.NewDatastoreFlag(ctx)
	cmd.DatastoreFlag.Register(ctx, f)

	f.StringVar(&cmd.adapterType, "a", string(types.VirtualDiskAdapterTypeLsiLogic), "Disk adapter")
	f.StringVar(&cmd.diskType, "d", "", "Disk format (thin, preallocated, ...); if empty, same as source")
	f.BoolVar(&cmd.force, "f", false, "Force")
}

func (cmd *cp) Process(ctx context.Context) error {
	return cmd.DatastoreFlag.Process(ctx)
}

func (cmd *cp) Usage() string {
	return "SRC DST"
}

func (cmd *cp) Description() string {
	return `Copy VMDK on DS, optionally converting format.

Examples:
  govc datastore.disk.copy disks/disk1.vmdk disks/disk2.vmdk
  govc datastore.disk.copy -d thin imported/stream.vmdk disks/thin.vmdk`
}

func (cmd *cp) Run(ctx context.Context, f *flag.FlagSet) error {
	if f.NArg() != 2 {
		return flag.ErrHelp
	}

	dc, err := cmd.Datacenter()
	if err != nil {
		return err
	}

	ds, err := cmd.Datastore()
	if err != nil {
		return err
	}

	m := object.NewVirtualDiskManager(ds.Client())

	src := ds.Path(f.Arg(0))
	dst := ds.Path(f.Arg(1))

	var destSpec types.BaseVirtualDiskSpec
	if cmd.diskType != "" {
		destSpec = &types.FileBackedVirtualDiskSpec{
			VirtualDiskSpec: types.VirtualDiskSpec{
				DiskType:    cmd.diskType,
				AdapterType: cmd.adapterType,
			},
		}
	}

	task, err := m.CopyVirtualDisk(ctx, src, dc, dst, dc, destSpec, cmd.force)
	if err != nil {
		return err
	}

	logger := cmd.ProgressLogger(fmt.Sprintf("Copying %s to %s...", src, dst))
	defer logger.Wait()

	_, err = task.WaitForResult(ctx, logger)
	return err
}
