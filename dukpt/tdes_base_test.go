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

// A.4.2 Initial Sequence test values from moov-io/dukpt.
const (
	pin           = "1234"
	pan           = "4012345678909"
	formatVersion = "ISO-0"
	data          = "4012345678909D987"
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
			Ksn:        pkg.HexDecode("FFFF9876543210E00001"),                             // //+YdlQyEOAAAQ==
			InitialKey: pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"),                 // asKS+qExW02FirOj19WTOg==
			CurrentKey: pkg.HexDecode("042666B49184CFA368DE9628D0397BC9"),                 // BCZmtJGEz6No3pYo0Dl7yQ==
			PinEnc:     pkg.HexDecode("1B9C1845EB993A7A"),                                 // G5wYReuZOno=
			DataReqEnc: pkg.HexDecode("FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F"), // /A1Tt+of2p7miq8ucNm5UGIpviqpk/BP
			DataResEnc: pkg.HexDecode("1FCC89AF66222F27B903898BB2BC8589CDBFDE5EC6AFCC25"), // H8yJr2YiLye5A4mLsryFic2/3l7Gr8wl
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

			pinEnc, err := EncryptPin(ck, pin, pan, formatVersion)
			require.NoError(t, err)
			desPinEnc, err := des.EncryptPin(desCk, pin, pan, formatVersion)
			require.NoError(t, err)
			require.Equal(t, pinEnc, desPinEnc)
			require.Equal(t, moov.PinEnc, pinEnc)

			decPin, err := DecryptPin(ck, pinEnc, pan, formatVersion)
			require.NoError(t, err)
			require.Equal(t, pin, decPin)

			reqEnc, err := EncryptData(ck, nil, data, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, moov.DataReqEnc, reqEnc)

			decReq, err := DecryptData(ck, reqEnc, nil, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, data, decReq[:len(data)])

			resEnc, err := EncryptData(ck, nil, data, pkg.ActionResponse)
			require.NoError(t, err)
			require.Equal(t, moov.DataResEnc, resEnc)

			decRes, err := DecryptData(ck, resEnc, nil, pkg.ActionResponse)
			require.NoError(t, err)
			require.Equal(t, data, decRes[:len(data)])

			// An empty IV must behave like the nil default instead of panicking.
			decEmptyIV, err := DecryptData(ck, reqEnc, []byte{}, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, decReq, decEmptyIV)
		})
	}
}
