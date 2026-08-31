// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.k6.io/k6/v2/js/modulestest"
)

func Test_module(t *testing.T) { //nolint:tparallel
	t.Parallel()

	runtime := modulestest.NewRuntime(t)

	err := runtime.SetupModuleSystem(map[string]any{ImportPath: new(RootModule)}, nil, nil)
	require.NoError(t, err)

	_, err = runtime.RunOnEventLoop(`let mod = require("` + ImportPath + `")`)
	require.NoError(t, err)

	tests := []struct {
		name  string
		check string
	}{
		{name: "getDesTcFromKsn(ksn)", check: `mod.getDesTcFromKsn(new Uint8Array([0xFF, 0xFF, 0x98, 0x76, 0x54, 0x32, 0x10, 0xE0, 0x00, 0x01])) === 1`},
		{name: "getDesTcFromKsnAsBase64(ksn)", check: `mod.getDesTcFromKsnAsBase64("//+YdlQyEOAAAQ==") === 1`},
		{name: "getDesTcFromKsnAsHex(ksn)", check: `mod.getDesTcFromKsnAsHex("FFFF9876543210E00001") === 1`},
		{name: "generateNextDesKsn(ksn)", check: `new Uint8Array(mod.generateNextDesKsn(new Uint8Array([0xFF, 0xFF, 0x98, 0x76, 0x54, 0x32, 0x10, 0xE0, 0x00, 0x01]))).toBase64() === "//+YdlQyEOAAAg=="`},
		{name: "generateNextDesKsn(ksn) keeps the input", check: `(() => { const ksn = Uint8Array.fromBase64("//+YdlQyEOAAAQ=="); mod.generateNextDesKsn(ksn); return ksn.toBase64() === "//+YdlQyEOAAAQ==" })()`},
		{name: "generateNextDesKsnAsBase64(ksn)", check: `mod.generateNextDesKsnAsBase64("//+YdlQyEOAAAQ==") === "//+YdlQyEOAAAg=="`},
		{name: "generateNextDesKsnAsHex(ksn)", check: `mod.generateNextDesKsnAsHex("FFFF9876543210E00001") === "FFFF9876543210E00002"`},
	}
	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			got, err := runtime.RunOnEventLoop(tt.check)

			require.NoError(t, err)
			require.True(t, got.ToBoolean())
		})
	}
}
