// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package crypto

import (
	"context"
	"flag"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/vmware/govmomi/cli"
	"github.com/vmware/govmomi/cli/flags"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type info struct {
	*flags.HostSystemFlag
	*flags.OutputFlag
}

func init() {
	cli.Register("host.crypto.info", &info{})
}

func (cmd *info) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.HostSystemFlag, ctx = flags.NewHostSystemFlag(ctx)
	cmd.HostSystemFlag.Register(ctx, f)
	cmd.OutputFlag, ctx = flags.NewOutputFlag(ctx)
	cmd.OutputFlag.Register(ctx, f)
}

func (cmd *info) Process(ctx context.Context) error {
	if err := cmd.HostSystemFlag.Process(ctx); err != nil {
		return err
	}
	return cmd.OutputFlag.Process(ctx)
}

func (cmd *info) Description() string {
	return `Display host crypto state.

CryptoState values: incapable, prepared, safe, pendingIncapable.

Examples:
  govc host.crypto.info
  govc host.crypto.info -host esxi01
  govc host.crypto.info -json`
}

type hostCryptoInfo struct {
	CryptoState string              `json:"cryptoState"`
	CryptoKeyId *types.CryptoKeyId `json:"cryptoKeyId,omitempty"`
}

type infoResult struct {
	hostCryptoInfo
}

func (r *infoResult) Write(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "%s\t%s\t%s\n", "CryptoState", "KeyId", "ProviderId")
	keyID, providerID := "-", "-"
	if r.CryptoKeyId != nil {
		keyID = r.CryptoKeyId.KeyId
		if r.CryptoKeyId.ProviderId != nil {
			providerID = r.CryptoKeyId.ProviderId.Id
		}
	}
	fmt.Fprintf(tw, "%s\t%s\t%s\n", r.CryptoState, keyID, providerID)
	return tw.Flush()
}

func (cmd *info) Run(ctx context.Context, f *flag.FlagSet) error {
	host, err := cmd.HostSystem()
	if err != nil {
		return err
	}

	var hs mo.HostSystem
	err = host.Properties(ctx, host.Reference(), []string{"runtime.cryptoState", "runtime.cryptoKeyId"}, &hs)
	if err != nil {
		return err
	}

	return cmd.WriteResult(&infoResult{hostCryptoInfo{
		CryptoState: hs.Runtime.CryptoState,
		CryptoKeyId: hs.Runtime.CryptoKeyId,
	}})
}
