// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"fmt"
	"testing"

	"github.com/moov-io/dukpt/pkg"
	"github.com/stretchr/testify/require"
)

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
			rawKsn := pkg.HexDecode(tt.ksn)

			next, err := GenerateNextDesKsn(rawKsn)
			require.NoError(t, err)
			require.Equal(t, pkg.HexDecode(tt.next), next)
			require.Equal(t, pkg.HexDecode(tt.ksn), rawKsn, "input must not be modified")

			desNext, err := pkg.GenerateNextDesKsn(pkg.HexDecode(tt.ksn))
			require.NoError(t, err)
			require.Equal(t, desNext, next)
		})
	}

	// Counter exhausted.
	_, err := GenerateNextDesKsn(pkg.HexDecode("FFFF9876543210FFF800"))
	require.Error(t, err)

	// Short KSN errors instead of panicking.
	_, err = GenerateNextDesKsn([]byte{0x01, 0x02})
	require.Error(t, err)
	_, err = GenerateNextDesKsn(nil)
	require.Error(t, err)
}
