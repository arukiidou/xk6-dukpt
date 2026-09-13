// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

// Package des is moov-io wrapper for des DUKPT (Derived Unique Key Per Transaction) cryptography operations.
package des

import (
	"maps"

	"go.k6.io/k6/v2/js/modules"
)

const ImportPath = "k6/x/dukpt/des"

type DesModule struct{}

func (*DesModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &module{vu}
}

type module struct {
	vu modules.VU
}

func (m *module) Exports() modules.Exports {
	base := map[string]any{
		"derivationOfInitialKey":      DerivationOfInitialKey,
		"deriveCurrentTransactionKey": DeriveCurrentTransactionKey,
		"encryptPin":                  EncryptPin,
		"decryptPin":                  DecryptPin,
		"encryptData":                 EncryptData,
		"decryptData":                 DecryptData,
		"generateMac":                 GenerateMac,
	}
	base64 := map[string]any{
		"derivationOfInitialKeyAsBase64":      DerivationOfInitialKeyAsBase64,
		"deriveCurrentTransactionKeyAsBase64": DeriveCurrentTransactionKeyAsBase64,
		"encryptPinAsBase64":                  EncryptPinAsBase64,
		"decryptPinAsBase64":                  DecryptPinAsBase64,
		"encryptDataAsBase64":                 EncryptDataAsBase64,
		"decryptDataAsBase64":                 DecryptDataAsBase64,
		"generateMacAsBase64":                 GenerateMacAsBase64,
	}
	maps.Copy(base, base64)

	return modules.Exports{
		Named: base,
	}
}

var _ modules.Module = (*DesModule)(nil)
