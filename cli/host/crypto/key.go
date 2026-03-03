// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package crypto

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// hexToKeyData converts a hex-encoded key to the base64-encoded string
// expected by the vSphere CryptoKeyPlain.KeyData field.
func hexToKeyData(hexKey string) (string, error) {
	raw, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", fmt.Errorf("-key-data: %w (expected hex-encoded key material)", err)
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}
