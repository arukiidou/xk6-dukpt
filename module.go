// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package xk6dukpt

import (
	"maps"

	"github.com/arukiidou/xk6-dukpt/dukpt"
	"go.k6.io/k6/v2/js/modules"
)

type rootModule struct{}

func (*rootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &module{vu}
}

type module struct {
	vu modules.VU
}

func (m *module) Exports() modules.Exports {
	std := map[string]any{
		"derivationOfInitialKey":      dukpt.DerivationOfInitialKey,
		"deriveCurrentTransactionKey": dukpt.DeriveCurrentTransactionKey,
		"encryptPin":                  dukpt.EncryptPin,
		"decryptPin":                  dukpt.DecryptPin,
		"encryptData":                 dukpt.EncryptData,
		"decryptData":                 dukpt.DecryptData,
		"generateMac":                 dukpt.GenerateMac,
	}
	base64 := map[string]any{
		"derivationOfInitialKeyAsBase64":      dukpt.DerivationOfInitialKeyAsBase64,
		"deriveCurrentTransactionKeyAsBase64": dukpt.DeriveCurrentTransactionKeyAsBase64,
		"encryptPinAsBase64":                  dukpt.EncryptPinAsBase64,
		"decryptPinAsBase64":                  dukpt.DecryptPinAsBase64,
		"encryptDataAsBase64":                 dukpt.EncryptDataAsBase64,
		"decryptDataAsBase64":                 dukpt.DecryptDataAsBase64,
		"generateMacAsBase64":                 dukpt.GenerateMacAsBase64,
	}
	maps.Copy(std, base64)

	return modules.Exports{
		Named: std,
	}
}

var _ modules.Module = (*rootModule)(nil)
