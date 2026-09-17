// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"fmt"
	"testing"

	"github.com/grafana/sobek"
	"github.com/moov-io/dukpt/pkg"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/v2/js/modulestest"
)

// newTestModule returns a module backed by a test runtime, plus an unwrap
// helper that flattens its ArrayBuffer results into raw bytes.
func newTestModule(t *testing.T) (*module, func(*sobek.ArrayBuffer, error) []byte) {
	t.Helper()

	instance, ok := new(RootModule).NewModuleInstance(modulestest.NewRuntime(t).VU).(*module)
	require.True(t, ok)

	unwrap := func(ab *sobek.ArrayBuffer, err error) []byte {
		t.Helper()

		require.NoError(t, err)
		require.NotNil(t, ab)
		return ab.Bytes()
	}
	return instance, unwrap
}

// A.4.2 Initial Sequence KSNs from moov-io/dukpt.
var moovInitialKsns = []string{
	"FFFF9876543210E00001", // //+YdlQyEOAAAQ==
	"FFFF9876543210E00002", // //+YdlQyEOAAAg==
	"FFFF9876543210E00003", // //+YdlQyEOAAAw==
	"FFFF9876543210E00004", // //+YdlQyEOAABA==
	"FFFF9876543210E00005", // //+YdlQyEOAABQ==
}

func TestGetDesTcFromKsn(t *testing.T) {
	t.Parallel()

	for index, ksn := range moovInitialKsns {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, ksn), func(t *testing.T) {
			t.Parallel()

			rawKsn := pkg.HexDecode(ksn)

			tc := GetDesTcFromKsn(rawKsn)
			require.Equal(t, pkg.GetDesTcFromKsn(rawKsn), tc)
			require.Equal(t, uint32(index+1), tc)
		})
	}

	// Low 21 bits only.
	require.Equal(t, uint32(0x1FFFFF), GetDesTcFromKsn(pkg.HexDecode("FFFF9876543210FFFFFF")))

	// KSN length is not validated.
	require.Equal(t, uint32(0), GetDesTcFromKsn([]byte{0x01, 0x02, 0x03}))
	require.Equal(t, uint32(0), GetDesTcFromKsn(nil))
}

func TestGenerateNextDesKsn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ksn  string
		next string
	}{
		{ksn: "FFFF9876543210E00001", next: "FFFF9876543210E00002"},
		{ksn: "FFFF9876543210E00004", next: "FFFF9876543210E00005"},
		// Counters with more than 10 bits set are skipped.
		{ksn: "FFFF9876543210E003FF", next: "FFFF9876543210E00400"},
		{ksn: "FFFF9876543210EFFC00", next: "FFFF9876543210F00000"},
	}
	for _, tt := range tests {
		t.Run(tt.ksn, func(t *testing.T) {
			t.Parallel()

			m, unwrap := newTestModule(t)
			rawKsn := pkg.HexDecode(tt.ksn)

			next := unwrap(m.GenerateNextDesKsn(rawKsn))
			require.Equal(t, pkg.HexDecode(tt.next), next)
			require.Equal(t, pkg.HexDecode(tt.ksn), rawKsn, "input must not be modified")

			desNext, err := pkg.GenerateNextDesKsn(pkg.HexDecode(tt.ksn))
			require.NoError(t, err)
			require.Equal(t, desNext, next)
		})
	}

	m, _ := newTestModule(t)

	// Counter exhausted.
	_, err := m.GenerateNextDesKsn(pkg.HexDecode("FFFF9876543210FFF800"))
	require.Error(t, err)

	// Short KSN errors instead of panicking.
	_, err = m.GenerateNextDesKsn([]byte{0x01, 0x02})
	require.Error(t, err)
	_, err = m.GenerateNextDesKsn(nil)
	require.Error(t, err)
}
