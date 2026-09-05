// Authors: arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package xk6dukpt

import (
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
	return modules.Exports{
		Named: map[string]any{
			"derivationOfInitialKeyAsBase64":      dukpt.DerivationOfInitialKeyAsBase64,
			"deriveCurrentTransactionKeyAsBase64": dukpt.DeriveCurrentTransactionKeyAsBase64,
		},
	}
}

var _ modules.Module = (*rootModule)(nil)
