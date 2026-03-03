// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package key

import (
	"context"
	"flag"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/vmware/govmomi/cli"
	"github.com/vmware/govmomi/cli/flags"
	"github.com/vmware/govmomi/crypto"
	"github.com/vmware/govmomi/vim25/types"
)

type ls struct {
	*flags.ClientFlag
	*flags.OutputFlag
	limit    int
	provider string
}

func init() {
	cli.Register("kms.key.ls", &ls{})
}

func (cmd *ls) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.ClientFlag, ctx = flags.NewClientFlag(ctx)
	cmd.ClientFlag.Register(ctx, f)
	cmd.OutputFlag, ctx = flags.NewOutputFlag(ctx)
	cmd.OutputFlag.Register(ctx, f)
	f.IntVar(&cmd.limit, "n", 0, "Maximum number of keys to return (0 = unlimited)")
	f.StringVar(&cmd.provider, "p", "", "Filter by provider ID")
}

func (cmd *ls) Process(ctx context.Context) error {
	if err := cmd.ClientFlag.Process(ctx); err != nil {
		return err
	}
	return cmd.OutputFlag.Process(ctx)
}

func (cmd *ls) Description() string {
	return `List crypto keys known to the vCenter CryptoManager.

Examples:
  govc kms.key.ls
  govc kms.key.ls -n 100
  govc kms.key.ls -p my-kp`
}

type lsResult struct {
	Keys []types.CryptoKeyId `json:"keys"`
}

func (r *lsResult) Write(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "%s\t%s\n", "KeyId", "ProviderId")
	for _, k := range r.Keys {
		pid := ""
		if k.ProviderId != nil {
			pid = k.ProviderId.Id
		}
		fmt.Fprintf(tw, "%s\t%s\n", k.KeyId, pid)
	}
	return tw.Flush()
}

func (cmd *ls) Run(ctx context.Context, f *flag.FlagSet) error {
	c, err := cmd.Client()
	if err != nil {
		return err
	}

	m, err := crypto.GetManagerKmip(c)
	if err != nil {
		return err
	}

	var limit *int32
	if cmd.limit > 0 {
		n := int32(cmd.limit)
		limit = &n
	}

	keys, err := m.ListKeys(ctx, limit)
	if err != nil {
		return err
	}

	if cmd.provider != "" {
		var filtered []types.CryptoKeyId
		for _, k := range keys {
			if k.ProviderId != nil && k.ProviderId.Id == cmd.provider {
				filtered = append(filtered, k)
			}
		}
		keys = filtered
	}

	return cmd.WriteResult(&lsResult{Keys: keys})
}
