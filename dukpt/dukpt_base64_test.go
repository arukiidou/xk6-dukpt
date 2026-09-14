// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/moov-io/dukpt/pkg"
	"github.com/stretchr/testify/require"
)

func TestBase64GetDesTcFromKsn(t *testing.T) {
	t.Parallel()

	for index, ksn := range moovInitialKsns {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, ksn), func(t *testing.T) {
			rawKsn := pkg.HexDecode(ksn)

			tc, err := GetDesTcFromKsnAsBase64(base64.StdEncoding.EncodeToString(rawKsn))
			require.NoError(t, err)
			require.Equal(t, pkg.GetDesTcFromKsn(rawKsn), tc)
			require.Equal(t, uint32(index+1), tc)
		})
	}

	_, err := GetDesTcFromKsnAsBase64("!!!")
	require.Error(t, err)
}

func TestBase64GenerateNextDesKsn(t *testing.T) {
	t.Parallel()

	for index, ksn := range moovInitialKsns[:len(moovInitialKsns)-1] {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, ksn), func(t *testing.T) {
			next, err := GenerateNextDesKsnAsBase64(base64.StdEncoding.EncodeToString(pkg.HexDecode(ksn)))
			require.NoError(t, err)
			require.Equal(t, base64.StdEncoding.EncodeToString(pkg.HexDecode(moovInitialKsns[index+1])), next)
		})
	}

	_, err := GenerateNextDesKsnAsBase64("!!!")
	require.Error(t, err)

	// Counter exhausted.
	_, err = GenerateNextDesKsnAsBase64(base64.StdEncoding.EncodeToString(pkg.HexDecode("FFFF9876543210FFF800")))
	require.Error(t, err)
}
