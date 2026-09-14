// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

// Package dukpt is moov-io wrapper for DUKPT (Derived Unique Key Per Transaction) utilities.
package dukpt

import (
	"maps"

	"go.k6.io/k6/v2/js/modules"
)

const ImportPath = "k6/x/dukpt"

type RootModule struct{}

func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &module{vu}
}

type module struct {
	vu modules.VU
}

func (m *module) Exports() modules.Exports {
	base := map[string]any{
		"getDesTcFromKsn":    GetDesTcFromKsn,
		"generateNextDesKsn": GenerateNextDesKsn,
	}
	base64 := map[string]any{
		"getDesTcFromKsnAsBase64":    GetDesTcFromKsnAsBase64,
		"generateNextDesKsnAsBase64": GenerateNextDesKsnAsBase64,
	}
	hexs := map[string]any{
		"getDesTcFromKsnAsHex":    GetDesTcFromKsnAsHex,
		"generateNextDesKsnAsHex": GenerateNextDesKsnAsHex,
	}
	maps.Copy(base, base64)
	maps.Copy(base, hexs)

	return modules.Exports{
		Named: base,
	}
}

var _ modules.Module = (*RootModule)(nil)
