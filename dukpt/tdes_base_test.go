// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"fmt"
	"testing"

	"github.com/moov-io/dukpt/pkg"
	"github.com/moov-io/dukpt/pkg/des"
	"github.com/stretchr/testify/require"
)

type SequenceItem struct {
	Bdk         []byte
	Ksn         []byte
	InitialKey  []byte
	CurrentKey  []byte
	PinEnc      []byte
	DataReqEnc  []byte
	DataResEnc  []byte
	RequestMac  []byte
	ResponseMac []byte
}

func TestMoovCompatibility(t *testing.T) {
	t.Parallel()

	bdk := pkg.HexDecode("0123456789ABCDEFFEDCBA9876543210") // ASNFZ4mrze/+3LqYdlQyEA==
	var InitialSequence = []SequenceItem{
		{
			Bdk:        bdk,
			Ksn:        pkg.HexDecode("FFFF9876543210E00001"),             // //+YdlQyEOAAAQ==
			InitialKey: pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"), // asKS+qExW02FirOj19WTOg==
			CurrentKey: pkg.HexDecode("042666B49184CFA368DE9628D0397BC9"), // BCZmtJGEz6No3pYo0Dl7yQ==
		},
	}

	for index, moov := range InitialSequence {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, pkg.HexEncode(moov.Ksn)), func(t *testing.T) {

			ik, err := DerivationOfInitialKey(moov.Bdk, moov.Ksn)
			require.NoError(t, err)
			desIk, err := des.DerivationOfInitialKey(moov.Bdk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, ik, desIk)

			ck, err := DeriveCurrentTransactionKey(ik, moov.Ksn)
			require.NoError(t, err)
			desCk, err := des.DeriveCurrentTransactionKey(desIk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, ck, desCk)
		})
	}
}
