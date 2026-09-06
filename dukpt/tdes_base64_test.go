// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"encoding/base64"
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

type SequenceBase64Item struct {
	Bdk         string
	Ksn         string
	InitialKey  string
	CurrentKey  string
	PinEnc      string
	DataReqEnc  string
	DataResEnc  string
	RequestMac  string
	ResponseMac string
}

func newSequenceBase64Item(seq SequenceItem) SequenceBase64Item {
	return SequenceBase64Item{
		Bdk:         base64.StdEncoding.EncodeToString(seq.Bdk),
		Ksn:         base64.StdEncoding.EncodeToString(seq.Ksn),
		InitialKey:  base64.StdEncoding.EncodeToString(seq.InitialKey),
		CurrentKey:  base64.StdEncoding.EncodeToString(seq.CurrentKey),
		PinEnc:      base64.StdEncoding.EncodeToString(seq.PinEnc),
		DataReqEnc:  base64.StdEncoding.EncodeToString(seq.DataReqEnc),
		DataResEnc:  base64.StdEncoding.EncodeToString(seq.DataResEnc),
		RequestMac:  base64.StdEncoding.EncodeToString(seq.RequestMac),
		ResponseMac: base64.StdEncoding.EncodeToString(seq.ResponseMac),
	}
}

func Test_module_DUKPT(t *testing.T) {
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
			b64 := newSequenceBase64Item(moov)

			ik, err := DerivationOfInitialKeyAsBase64(b64.Bdk, b64.Ksn)
			require.NoError(t, err)
			require.Len(t, ik, 24)
			require.Equal(t, ik, base64.StdEncoding.EncodeToString(moov.InitialKey))

			ck, err := DeriveCurrentTransactionKeyAsBase64(ik, b64.Ksn)
			require.NoError(t, err)
			require.Len(t, ck, 24)
			require.Equal(t, ck, base64.StdEncoding.EncodeToString(moov.CurrentKey))
		})
	}
}

// Check moov-io/dukpt compatibility.
func TestBase64MoovCompatibility(t *testing.T) {
	t.Parallel()

	var InitialSequence = []SequenceItem{
		{
			Bdk:        pkg.HexDecode("0123456789ABCDEFFEDCBA9876543210"), // ASNFZ4mrze/+3LqYdlQyEA==
			Ksn:        pkg.HexDecode("FFFF9876543210E00001"),             // //+YdlQyEOAAAQ==
			InitialKey: pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"), // asKS+qExW02FirOj19WTOg==
			CurrentKey: pkg.HexDecode("042666B49184CFA368DE9628D0397BC9"), // BCZmtJGEz6No3pYo0Dl7yQ==
		}}

	for index, moov := range InitialSequence {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, pkg.HexEncode(moov.Ksn)), func(t *testing.T) {
			b64 := newSequenceBase64Item(moov)

			ik, _ := DerivationOfInitialKeyAsBase64(b64.Bdk, b64.Ksn)
			desIk, _ := des.DerivationOfInitialKey(moov.Bdk, moov.Ksn)
			require.Equal(t, ik, base64.StdEncoding.EncodeToString(desIk))

			ck, _ := DeriveCurrentTransactionKeyAsBase64(ik, b64.Ksn)
			desCk, _ := des.DeriveCurrentTransactionKey(desIk, moov.Ksn)
			require.Equal(t, ck, base64.StdEncoding.EncodeToString(desCk))
		})
	}
}
