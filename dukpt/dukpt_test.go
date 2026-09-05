// Authors: arukiidou <arukiidou@yahoo.co.jp>
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
	Ksn         []byte
	InitialKey  []byte
	CurrentKey  []byte
	PinEnc      []byte
	DataReqEnc  []byte
	DataResEnc  []byte
	RequestMac  []byte
	ResponseMac []byte
}

func Test_module_DUKPT(t *testing.T) {
	t.Parallel()

	bdk := pkg.HexDecode("0123456789ABCDEFFEDCBA9876543210") // ASNFZ4mrze/+3LqYdlQyEA==
	var InitialSequence = []SequenceItem{
		{
			Ksn:        pkg.HexDecode("FFFF9876543210E00001"),             // //+YdlQyEOAAAQ==
			InitialKey: pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"), // asKS+qExW02FirOj19WTOg==
			CurrentKey: pkg.HexDecode("042666B49184CFA368DE9628D0397BC9"), // BCZmtJGEz6No3pYo0Dl7yQ==
		}}

	b64 := base64.StdEncoding
	for index, item := range InitialSequence {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, pkg.HexEncode(item.Ksn)), func(t *testing.T) {
			ik, err := DerivationOfInitialKeyAsBase64(b64.EncodeToString(bdk), b64.EncodeToString(item.Ksn))
			require.NoError(t, err)
			require.Len(t, ik, 24)
			require.Equal(t, ik, b64.EncodeToString(item.InitialKey))

			ck, err := DeriveCurrentTransactionKeyAsBase64(ik, b64.EncodeToString(item.Ksn))
			require.NoError(t, err)
			require.Len(t, ck, 24)
			require.Equal(t, ck, b64.EncodeToString(item.CurrentKey))
		})
	}
}

// Check moov-io/dukpt compatibility.
func TestModuleMoovCompatibility(t *testing.T) {
	t.Parallel()

	bdkHex := pkg.HexDecode("0123456789ABCDEFFEDCBA9876543210") // ASNFZ4mrze/+3LqYdlQyEA==
	var InitialSequence = []SequenceItem{
		{
			Ksn:        pkg.HexDecode("FFFF9876543210E00001"),             // //+YdlQyEOAAAQ==
			InitialKey: pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"), // asKS+qExW02FirOj19WTOg==
			CurrentKey: pkg.HexDecode("042666B49184CFA368DE9628D0397BC9"), // BCZmtJGEz6No3pYo0Dl7yQ==
		}}

	b64 := base64.StdEncoding
	for index, item := range InitialSequence {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, pkg.HexEncode(item.Ksn)), func(t *testing.T) {
			bdk := b64.EncodeToString(bdkHex)
			ksn := b64.EncodeToString(item.Ksn)

			ik, _ := DerivationOfInitialKeyAsBase64(bdk, ksn)
			desIk, _ := des.DerivationOfInitialKey(bdkHex, item.Ksn)
			require.Equal(t, ik, b64.EncodeToString(desIk))

			ck, _ := DeriveCurrentTransactionKeyAsBase64(ik, ksn)
			desCk, _ := des.DeriveCurrentTransactionKey(desIk, item.Ksn)
			require.Equal(t, ck, b64.EncodeToString(desCk))
		})
	}
}
