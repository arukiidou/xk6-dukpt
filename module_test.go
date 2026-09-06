// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package xk6dukpt

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"go.k6.io/k6/v2/js/modulestest"
)

func Test_module(t *testing.T) { //nolint:tparallel
	t.Parallel()

	runtime := modulestest.NewRuntime(t)

	err := runtime.SetupModuleSystem(map[string]any{importPath: new(rootModule)}, nil, nil)
	require.NoError(t, err)

	_, err = runtime.RunOnEventLoop(`let mod = require("` + importPath + `")`)
	require.NoError(t, err)

	tests := []struct {
		name  string
		check string
	}{
		{name: "derivationOfInitialKeyAsBase64(bdk, ksn)", check: `mod.derivationOfInitialKeyAsBase64("ASNFZ4mrze/+3LqYdlQyEA==", "//+YdlQyEOAAAQ==") === "asKS+qExW02FirOj19WTOg=="`},
		{name: "deriveCurrentTransactionKeyAsBase64(ik, ksn)", check: `mod.deriveCurrentTransactionKeyAsBase64("asKS+qExW02FirOj19WTOg==", "//+YdlQyEOAAAQ==") === "BCZmtJGEz6No3pYo0Dl7yQ=="`},
	}
	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			got, err := runtime.RunOnEventLoop(tt.check)

			require.NoError(t, err)
			require.True(t, got.ToBoolean())
		})
	}
}
