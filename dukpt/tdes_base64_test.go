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

// Check moov-io/dukpt compatibility.
func TestBase64MoovCompatibility(t *testing.T) {
	t.Parallel()

	for index, moov := range moovInitialSequence {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, pkg.HexEncode(moov.Ksn)), func(t *testing.T) {
			b64 := newSequenceBase64Item(moov)

			ik, err := DerivationOfInitialKeyAsBase64(b64.Bdk, b64.Ksn)
			require.NoError(t, err)
			desIk, err := des.DerivationOfInitialKey(moov.Bdk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, ik, base64.StdEncoding.EncodeToString(desIk))
			require.Equal(t, b64.InitialKey, ik)

			ck, err := DeriveCurrentTransactionKeyAsBase64(ik, b64.Ksn)
			require.NoError(t, err)
			desCk, err := des.DeriveCurrentTransactionKey(desIk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, ck, base64.StdEncoding.EncodeToString(desCk))
			require.Equal(t, b64.CurrentKey, ck)

			pinEnc, err := EncryptPinAsBase64(ck, pin, pan, formatVersion)
			require.NoError(t, err)
			desPinEnc, err := des.EncryptPin(desCk, pin, pan, formatVersion)
			require.NoError(t, err)
			require.Equal(t, pinEnc, base64.StdEncoding.EncodeToString(desPinEnc))
			require.Equal(t, b64.PinEnc, pinEnc)

			decPin, err := DecryptPinAsBase64(ck, pinEnc, pan, formatVersion)
			require.NoError(t, err)
			require.Equal(t, pin, decPin)

			reqEnc, err := EncryptDataAsBase64(ck, "", data, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, b64.DataReqEnc, reqEnc)

			decReq, err := DecryptDataAsBase64(ck, reqEnc, "", pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, data, decReq[:len(data)])

			resEnc, err := EncryptDataAsBase64(ck, "", data, pkg.ActionResponse)
			require.NoError(t, err)
			require.Equal(t, b64.DataResEnc, resEnc)

			decRes, err := DecryptDataAsBase64(ck, resEnc, "", pkg.ActionResponse)
			require.NoError(t, err)
			require.Equal(t, data, decRes[:len(data)])
		})
	}
}
