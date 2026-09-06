// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package xk6dukpt

import (
	"maps"

	"github.com/arukiidou/xk6-dukpt/dukpt"
	"github.com/moov-io/dukpt/pkg/des"
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
		"derivationOfInitialKey":      des.DerivationOfInitialKey,
		"deriveCurrentTransactionKey": des.DeriveCurrentTransactionKey,
	}
	base64 := map[string]any{
		"derivationOfInitialKeyAsBase64":      dukpt.DerivationOfInitialKeyAsBase64,
		"deriveCurrentTransactionKeyAsBase64": dukpt.DeriveCurrentTransactionKeyAsBase64,
	}
	maps.Copy(std, base64)

	return modules.Exports{
		Named: std,
	}
}

var _ modules.Module = (*rootModule)(nil)
