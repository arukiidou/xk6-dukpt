// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package des

import (
	"fmt"
	"strings"
	"testing"

	"github.com/moov-io/dukpt/pkg"
	"github.com/moov-io/dukpt/pkg/des"
	"github.com/stretchr/testify/require"
)

type SequenceHexItem struct {
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

func newSequenceHexItem(seq SequenceItem) SequenceHexItem {
	return SequenceHexItem{
		Bdk:         encodeHex(seq.Bdk),
		Ksn:         encodeHex(seq.Ksn),
		InitialKey:  encodeHex(seq.InitialKey),
		CurrentKey:  encodeHex(seq.CurrentKey),
		PinEnc:      encodeHex(seq.PinEnc),
		DataReqEnc:  encodeHex(seq.DataReqEnc),
		DataResEnc:  encodeHex(seq.DataResEnc),
		RequestMac:  encodeHex(seq.RequestMac),
		ResponseMac: encodeHex(seq.ResponseMac),
	}
}

// Check moov-io/dukpt compatibility.
func TestHexMoovCompatibility(t *testing.T) {
	t.Parallel()

	for index, moov := range moovInitialSequence {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, pkg.HexEncode(moov.Ksn)), func(t *testing.T) {
			hx := newSequenceHexItem(moov)

			ik, err := DerivationOfInitialKeyAsHex(hx.Bdk, hx.Ksn)
			require.NoError(t, err)
			desIk, err := des.DerivationOfInitialKey(moov.Bdk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, ik, encodeHex(desIk))
			require.Equal(t, hx.InitialKey, ik)

			ck, err := DeriveCurrentTransactionKeyAsHex(ik, hx.Ksn)
			require.NoError(t, err)
			desCk, err := des.DeriveCurrentTransactionKey(desIk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, ck, encodeHex(desCk))
			require.Equal(t, hx.CurrentKey, ck)

			pinEnc, err := EncryptPinAsHex(ck, pin, pan, formatVersion)
			require.NoError(t, err)
			desPinEnc, err := des.EncryptPin(desCk, pin, pan, formatVersion)
			require.NoError(t, err)
			require.Equal(t, pinEnc, encodeHex(desPinEnc))
			require.Equal(t, hx.PinEnc, pinEnc)

			decPin, err := DecryptPinAsHex(ck, pinEnc, pan, formatVersion)
			require.NoError(t, err)
			require.Equal(t, pin, decPin)

			reqEnc, err := EncryptDataAsHex(ck, "", data, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, hx.DataReqEnc, reqEnc)

			decReq, err := DecryptDataAsHex(ck, reqEnc, "", pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, data, decReq[:len(data)])

			resEnc, err := EncryptDataAsHex(ck, "", data, pkg.ActionResponse)
			require.NoError(t, err)
			require.Equal(t, hx.DataResEnc, resEnc)

			decRes, err := DecryptDataAsHex(ck, resEnc, "", pkg.ActionResponse)
			require.NoError(t, err)
			require.Equal(t, data, decRes[:len(data)])

			reqMac, err := GenerateMacAsHex(ck, data, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, hx.RequestMac, reqMac)

			resMac, err := GenerateMacAsHex(ck, data, pkg.ActionResponse)
			require.NoError(t, err)
			require.Equal(t, hx.ResponseMac, resMac)
		})
	}
}

func TestHexAcceptsLowercase(t *testing.T) {
	t.Parallel()

	hx := newSequenceHexItem(moovInitialSequence[0])

	ik, err := DerivationOfInitialKeyAsHex(strings.ToLower(hx.Bdk), strings.ToLower(hx.Ksn))
	require.NoError(t, err)
	require.Equal(t, hx.InitialKey, ik)
}

func TestHexInvalidInput(t *testing.T) {
	t.Parallel()

	hx := newSequenceHexItem(moovInitialSequence[0])

	tests := []struct {
		name string
		call func() error
	}{
		{name: "invalid character", call: func() error {
			_, err := DerivationOfInitialKeyAsHex("ZZ", hx.Ksn)
			return err
		}},
		{name: "odd length", call: func() error {
			_, err := DeriveCurrentTransactionKeyAsHex(hx.InitialKey, "ABC")
			return err
		}},
		{name: "invalid iv", call: func() error {
			_, err := EncryptDataAsHex(hx.CurrentKey, "XY", data, pkg.ActionRequest)
			return err
		}},
		{name: "invalid ciphertext", call: func() error {
			_, err := DecryptPinAsHex(hx.CurrentKey, "G5wYReuZOno=", pan, formatVersion)
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Error(t, tt.call())
		})
	}
}
