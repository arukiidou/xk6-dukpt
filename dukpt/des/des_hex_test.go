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
		Bdk:         encodeUpperHex(seq.Bdk),
		Ksn:         encodeUpperHex(seq.Ksn),
		InitialKey:  encodeUpperHex(seq.InitialKey),
		CurrentKey:  encodeUpperHex(seq.CurrentKey),
		PinEnc:      encodeUpperHex(seq.PinEnc),
		DataReqEnc:  encodeUpperHex(seq.DataReqEnc),
		DataResEnc:  encodeUpperHex(seq.DataResEnc),
		RequestMac:  encodeUpperHex(seq.RequestMac),
		ResponseMac: encodeUpperHex(seq.ResponseMac),
	}
}

// Check moov-io/dukpt compatibility.
func TestHexMoovCompatibility(t *testing.T) {
	t.Parallel()

	for index, moov := range moovInitialSequence {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, pkg.HexEncode(moov.Ksn)), func(t *testing.T) {
			t.Parallel()

			hx := newSequenceHexItem(moov)

			ik, err := DerivationOfInitialKeyAsHex(hx.Bdk, hx.Ksn)
			require.NoError(t, err)
			desIk, err := des.DerivationOfInitialKey(moov.Bdk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, encodeUpperHex(desIk), ik)
			require.Equal(t, hx.InitialKey, ik)

			ck, err := DeriveCurrentTransactionKeyAsHex(ik, hx.Ksn)
			require.NoError(t, err)
			desCk, err := des.DeriveCurrentTransactionKey(desIk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, encodeUpperHex(desCk), ck)
			require.Equal(t, hx.CurrentKey, ck)

			pinEnc, err := EncryptPinAsHex(ck, pin, pan, formatVersion)
			require.NoError(t, err)
			desPinEnc, err := des.EncryptPin(desCk, pin, pan, formatVersion)
			require.NoError(t, err)
			require.Equal(t, encodeUpperHex(desPinEnc), pinEnc)
			require.Equal(t, hx.PinEnc, pinEnc)

			decPin, err := DecryptPinAsHex(ck, pinEnc, pan, formatVersion)
			require.NoError(t, err)
			require.Equal(t, pin, decPin)

			reqEnc, err := EncryptDataAsHex(ck, "", data, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, hx.DataReqEnc, reqEnc)

			decReq, err := DecryptDataAsHex(ck, reqEnc, "", pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, encodeUpperHex([]byte(data)), decReq[:len(data)*2])

			resEnc, err := EncryptDataAsHex(ck, "", data, pkg.ActionResponse)
			require.NoError(t, err)
			require.Equal(t, hx.DataResEnc, resEnc)

			decRes, err := DecryptDataAsHex(ck, resEnc, "", pkg.ActionResponse)
			require.NoError(t, err)
			require.Equal(t, encodeUpperHex([]byte(data)), decRes[:len(data)*2])

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

	// Every string argument is hex decoded, so each one has to reject garbage.
	tests := []struct {
		name string
		call func() error
	}{
		{name: "invalid bdk", call: func() error {
			_, err := DerivationOfInitialKeyAsHex("ZZ", hx.Ksn)
			return err
		}},
		{name: "invalid ksn for derivation", call: func() error {
			_, err := DerivationOfInitialKeyAsHex(hx.Bdk, "ZZ")
			return err
		}},
		{name: "invalid ik", call: func() error {
			_, err := DeriveCurrentTransactionKeyAsHex("ZZ", hx.Ksn)
			return err
		}},
		{name: "odd length ksn", call: func() error {
			_, err := DeriveCurrentTransactionKeyAsHex(hx.InitialKey, "ABC")
			return err
		}},
		{name: "invalid current key for encrypt pin", call: func() error {
			_, err := EncryptPinAsHex("ZZ", pin, pan, formatVersion)
			return err
		}},
		{name: "invalid current key for decrypt pin", call: func() error {
			_, err := DecryptPinAsHex("ZZ", hx.PinEnc, pan, formatVersion)
			return err
		}},
		{name: "invalid pin ciphertext", call: func() error {
			// The base64 form of the same pin block is not valid hex.
			_, err := DecryptPinAsHex(hx.CurrentKey, "G5wYReuZOno=", pan, formatVersion)
			return err
		}},
		{name: "invalid current key for encrypt data", call: func() error {
			_, err := EncryptDataAsHex("ZZ", "", data, pkg.ActionRequest)
			return err
		}},
		{name: "invalid iv character for encrypt data", call: func() error {
			_, err := EncryptDataAsHex(hx.CurrentKey, "XY", data, pkg.ActionRequest)
			return err
		}},
		{name: "invalid current key for decrypt data", call: func() error {
			_, err := DecryptDataAsHex("ZZ", hx.DataReqEnc, "", pkg.ActionRequest)
			return err
		}},
		{name: "invalid data ciphertext", call: func() error {
			_, err := DecryptDataAsHex(hx.CurrentKey, "ZZ", "", pkg.ActionRequest)
			return err
		}},
		{name: "invalid iv character for decrypt data", call: func() error {
			_, err := DecryptDataAsHex(hx.CurrentKey, hx.DataReqEnc, "XY", pkg.ActionRequest)
			return err
		}},
		{name: "invalid current key for generate mac", call: func() error {
			_, err := GenerateMacAsHex("ZZ", data, pkg.ActionRequest)
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

// Errors moov-io raises on well formed hex have to reach the caller too.
func TestHexMoovErrors(t *testing.T) {
	t.Parallel()

	hx := newSequenceHexItem(moovInitialSequence[0])

	tests := []struct {
		name string
		call func() error
	}{
		{name: "short bdk", call: func() error {
			_, err := DerivationOfInitialKeyAsHex(hx.Bdk[:len(hx.Bdk)-2], hx.Ksn)
			return err
		}},
		{name: "unknown pin format for encrypt", call: func() error {
			_, err := EncryptPinAsHex(hx.CurrentKey, pin, pan, "ISO-9")
			return err
		}},
		{name: "unknown pin format for decrypt", call: func() error {
			_, err := DecryptPinAsHex(hx.CurrentKey, hx.PinEnc, pan, "ISO-9")
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

// The hex form of [TestIVLength]: an empty IV is the zero vector and anything
// that is not one block is rejected before it reaches moov-io.
func TestHexIVLength(t *testing.T) {
	t.Parallel()

	hx := newSequenceHexItem(moovInitialSequence[0])

	tests := []struct {
		name string
		iv   string
		// block is the one-block IV that iv has to behave like, empty when iv is rejected.
		block string
	}{
		{name: "empty", iv: "", block: "0000000000000000"},
		{name: "zero", iv: "0000000000000000", block: "0000000000000000"},
		{name: "exact", iv: "0102030405060708", block: "0102030405060708"},
		{name: "lowercase", iv: "0a0b0c0d0e0f1011", block: "0A0B0C0D0E0F1011"},
		{name: "short", iv: "0102"},
		{name: "long", iv: "0102030405060708090A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.block == "" {
				wantErr := fmt.Sprintf("iv must be %d bytes, got %d", desBlockLen, len(tt.iv)/2)

				_, err := EncryptDataAsHex(hx.CurrentKey, tt.iv, data, pkg.ActionRequest)
				require.ErrorContains(t, err, wantErr)

				_, err = DecryptDataAsHex(hx.CurrentKey, hx.DataReqEnc, tt.iv, pkg.ActionRequest)
				require.ErrorContains(t, err, wantErr)
				return
			}

			ciphertext, err := EncryptDataAsHex(hx.CurrentKey, tt.iv, data, pkg.ActionRequest)
			require.NoError(t, err)
			withBlock, err := EncryptDataAsHex(hx.CurrentKey, tt.block, data, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, withBlock, ciphertext)

			plaintext, err := DecryptDataAsHex(hx.CurrentKey, ciphertext, tt.iv, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, encodeUpperHex([]byte(data)), plaintext[:len(data)*2])
		})
	}
}
