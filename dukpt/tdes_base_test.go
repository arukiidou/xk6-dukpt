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

// Common BDK for all test sequences, used to derive the initial key.
var moovTestBdk = pkg.HexDecode("0123456789ABCDEFFEDCBA9876543210") // ASNFZ4mrze/+3LqYdlQyEA==

// The A.4.2 Initial Sequence vectors, shared by the raw and the base64 suites.
// Every sequence derives from the same BDK, so the initial key is constant and
// only the KSN counter advances.
//
// RequestMac and ResponseMac stay empty: this package exposes no MAC wrapper.
var moovInitialSequence = []SequenceItem{
	{
		Bdk:        moovTestBdk,
		Ksn:        pkg.HexDecode("FFFF9876543210E00001"),                             // //+YdlQyEOAAAQ==
		InitialKey: pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"),                 // asKS+qExW02FirOj19WTOg==
		CurrentKey: pkg.HexDecode("042666B49184CFA368DE9628D0397BC9"),                 // BCZmtJGEz6No3pYo0Dl7yQ==
		PinEnc:     pkg.HexDecode("1B9C1845EB993A7A"),                                 // G5wYReuZOno=
		DataReqEnc: pkg.HexDecode("FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F"), // /A1Tt+of2p7miq8ucNm5UGIpviqpk/BP
		DataResEnc: pkg.HexDecode("1FCC89AF66222F27B903898BB2BC8589CDBFDE5EC6AFCC25"), // H8yJr2YiLye5A4mLsryFic2/3l7Gr8wl
	},
	{
		Bdk:        moovTestBdk,
		Ksn:        pkg.HexDecode("FFFF9876543210E00002"),                             // //+YdlQyEOAAAg==
		InitialKey: pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"),                 // asKS+qExW02FirOj19WTOg==
		CurrentKey: pkg.HexDecode("C46551CEF9FD24B0AA9AD834130D3BC7"),                 // xGVRzvn9JLCqmtg0Ew07xw==
		PinEnc:     pkg.HexDecode("10A01C8D02C69107"),                                 // EKAcjQLGkQc=
		DataReqEnc: pkg.HexDecode("A2B4E70F846E63D68775B7215EB4563DFD3037244C61CC13"), // orTnD4RuY9aHdbchXrRWPf0wNyRMYcwT
		DataResEnc: pkg.HexDecode("5B692A6B1FDD5E25B0DEFAFDE1672E402F8011360CFF3508"), // W2kqax/dXiWw3vr94WcuQC+AETYM/zUI
	},
	{
		Bdk:        moovTestBdk,
		Ksn:        pkg.HexDecode("FFFF9876543210E00003"),                             // //+YdlQyEOAAAw==
		InitialKey: pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"),                 // asKS+qExW02FirOj19WTOg==
		CurrentKey: pkg.HexDecode("0DF3D9422ACA56E547676D07AD6BADFA"),                 // DfPZQirKVuVHZ20HrWut+g==
		PinEnc:     pkg.HexDecode("18DC07B94797B466"),                                 // GNwHuUeXtGY=
		DataReqEnc: pkg.HexDecode("BD751E65F10E75B6C1D5B1D283496A36C2DE21D993C387A7"), // vXUeZfEOdbbB1bHSg0lqNsLeIdmTw4en
		DataResEnc: pkg.HexDecode("345992D4163E4926C927BFD8ABD5D76F087A9CE81D5A27B8"), // NFmS1BY+SSbJJ7/Yq9XXbwh6nOgdWie4
	},
	{
		Bdk:        moovTestBdk,
		Ksn:        pkg.HexDecode("FFFF9876543210E00004"),                             // //+YdlQyEOAABA==
		InitialKey: pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"),                 // asKS+qExW02FirOj19WTOg==
		CurrentKey: pkg.HexDecode("279C0F6AEED0BE652B2C733E1383AE91"),                 // J5wPau7QvmUrLHM+E4OukQ==
		PinEnc:     pkg.HexDecode("0BC79509D5645DF7"),                                 // C8eVCdVkXfc=
		DataReqEnc: pkg.HexDecode("1118F50947441BBDA3C8C70220021A12EC31CC473F7215F4"), // ERj1CUdEG72jyMcCIAIaEuwxzEc/chX0
		DataResEnc: pkg.HexDecode("418C7413576C0D1819E785D3807AF32334231FDEC23414DB"), // QYx0E1dsDRgZ54XTgHrzIzQjH97CNBTb
	},
	{
		Bdk:        moovTestBdk,
		Ksn:        pkg.HexDecode("FFFF9876543210E00005"),                             // //+YdlQyEOAABQ==
		InitialKey: pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"),                 // asKS+qExW02FirOj19WTOg==
		CurrentKey: pkg.HexDecode("5F8DC6D2C845C125508DDC048093B83F"),                 // X43G0shFwSVQjdwEgJO4Pw==
		PinEnc:     pkg.HexDecode("5BC0AF22AD87B327"),                                 // W8CvIq2Hsyc=
		DataReqEnc: pkg.HexDecode("9FD7BD1EC28845ACA93367A9DA9317BD555C6B33AE22D365"), // n9e9HsKIRaypM2ep2pMXvVVcazOuItNl
		DataResEnc: pkg.HexDecode("7D4C109E49E83355A556AE949EED359F4404E7A2F0167C00"), // fUwQnknoM1WlVq6Unu01n0QE56LwFnwA
	},
}

func TestMoovCompatibility(t *testing.T) {
	t.Parallel()

	for index, moov := range moovInitialSequence {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, pkg.HexEncode(moov.Ksn)), func(t *testing.T) {

			ik, err := DerivationOfInitialKey(moov.Bdk, moov.Ksn)
			require.NoError(t, err)
			desIk, err := des.DerivationOfInitialKey(moov.Bdk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, ik, desIk)
			require.Equal(t, moov.InitialKey, ik)

			ck, err := DeriveCurrentTransactionKey(ik, moov.Ksn)
			require.NoError(t, err)
			desCk, err := des.DeriveCurrentTransactionKey(desIk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, ck, desCk)
			require.Equal(t, moov.CurrentKey, ck)

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
