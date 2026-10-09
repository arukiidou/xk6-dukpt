// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/moov-io/dukpt/pkg"
	"github.com/stretchr/testify/require"
)

func TestHexGetDesTcFromKsn(t *testing.T) {
	t.Parallel()

	for index, ksn := range moovInitialKsns {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, ksn), func(t *testing.T) {
			t.Parallel()

			tc, err := GetDesTcFromKsnAsHex(ksn)
			require.NoError(t, err)
			require.Equal(t, pkg.GetDesTcFromKsn(pkg.HexDecode(ksn)), tc)
			require.Equal(t, uint32(index+1), tc)
		})
	}
}

func TestHexAcceptsLowercase(t *testing.T) {
	t.Parallel()

	tc, err := GetDesTcFromKsnAsHex(strings.ToLower(moovInitialKsns[0]))
	require.NoError(t, err)
	require.Equal(t, uint32(1), tc)
}

func TestHexInvalidInput(t *testing.T) {
	t.Parallel()

	_, err := GetDesTcFromKsnAsHex("ZZ")
	require.Error(t, err)

	_, err = GetDesTcFromKsnAsHex("ABC")
	require.Error(t, err)

	_, err = GenerateNextDesKsnAsHex("ZZ")
	require.Error(t, err)

	// Below 4 bytes moov-io reads counter 0, which is an error here.
	_, err = GetDesTcFromKsnAsHex("0000FF")
	require.Error(t, err)

	// 3 bytes would read counter 0 and silently restart it at 1.
	_, err = GenerateNextDesKsnAsHex("E00008")
	require.Error(t, err)
}

func TestHexGenerateNextDesKsn(t *testing.T) {
	t.Parallel()

	for index, ksn := range moovInitialKsns[:len(moovInitialKsns)-1] {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, ksn), func(t *testing.T) {
			t.Parallel()

			next, err := GenerateNextDesKsnAsHex(ksn)
			require.NoError(t, err)
			require.Equal(t, moovInitialKsns[index+1], next)
		})
	}

	// Lowercase input, uppercase output.
	next, err := GenerateNextDesKsnAsHex(strings.ToLower(moovInitialKsns[0]))
	require.NoError(t, err)
	require.Equal(t, moovInitialKsns[1], next)

	// Counter exhausted.
	_, err = GenerateNextDesKsnAsHex("FFFF9876543210FFF800")
	require.Error(t, err)
}
