// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package des

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
			t.Parallel()

			b64 := newSequenceBase64Item(moov)

			ik, err := DerivationOfInitialKeyAsBase64(b64.Bdk, b64.Ksn)
			require.NoError(t, err)
			desIk, err := des.DerivationOfInitialKey(moov.Bdk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, base64.StdEncoding.EncodeToString(desIk), ik)
			require.Equal(t, b64.InitialKey, ik)

			ck, err := DeriveCurrentTransactionKeyAsBase64(ik, b64.Ksn)
			require.NoError(t, err)
			desCk, err := des.DeriveCurrentTransactionKey(desIk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, base64.StdEncoding.EncodeToString(desCk), ck)
			require.Equal(t, b64.CurrentKey, ck)

			pinEnc, err := EncryptPinAsBase64(ck, pin, pan, formatVersion)
			require.NoError(t, err)
			desPinEnc, err := des.EncryptPin(desCk, pin, pan, formatVersion)
			require.NoError(t, err)
			require.Equal(t, base64.StdEncoding.EncodeToString(desPinEnc), pinEnc)
			require.Equal(t, b64.PinEnc, pinEnc)

			decPin, err := DecryptPinAsBase64(ck, pinEnc, pan, formatVersion)
			require.NoError(t, err)
			require.Equal(t, pin, decPin)

			reqEnc, err := EncryptDataAsBase64(ck, "", data, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, b64.DataReqEnc, reqEnc)

			decReq, err := DecryptDataAsBase64(ck, reqEnc, "", pkg.ActionRequest)
			require.NoError(t, err)
			rawDecReq, err := base64.StdEncoding.DecodeString(decReq)
			require.NoError(t, err)
			require.Equal(t, data, string(rawDecReq[:len(data)]))

			resEnc, err := EncryptDataAsBase64(ck, "", data, pkg.ActionResponse)
			require.NoError(t, err)
			require.Equal(t, b64.DataResEnc, resEnc)

			decRes, err := DecryptDataAsBase64(ck, resEnc, "", pkg.ActionResponse)
			require.NoError(t, err)
			rawDecRes, err := base64.StdEncoding.DecodeString(decRes)
			require.NoError(t, err)
			require.Equal(t, data, string(rawDecRes[:len(data)]))

			reqMac, err := GenerateMacAsBase64(ck, data, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, b64.RequestMac, reqMac)

			resMac, err := GenerateMacAsBase64(ck, data, pkg.ActionResponse)
			require.NoError(t, err)
			require.Equal(t, b64.ResponseMac, resMac)
		})
	}
}

// The base64 form of [TestHexInvalidInput]: every string argument is decoded,
// so each one has to reject garbage.
func TestBase64InvalidInput(t *testing.T) {
	t.Parallel()

	b64 := newSequenceBase64Item(moovInitialSequence[0])

	// Not valid standard base64: "!" is outside the alphabet.
	const invalid = "!!!"

	tests := []struct {
		name string
		call func() error
	}{
		{name: "invalid bdk", call: func() error {
			_, err := DerivationOfInitialKeyAsBase64(invalid, b64.Ksn)
			return err
		}},
		{name: "invalid ksn for derivation", call: func() error {
			_, err := DerivationOfInitialKeyAsBase64(b64.Bdk, invalid)
			return err
		}},
		{name: "invalid ik", call: func() error {
			_, err := DeriveCurrentTransactionKeyAsBase64(invalid, b64.Ksn)
			return err
		}},
		{name: "invalid ksn for derive", call: func() error {
			_, err := DeriveCurrentTransactionKeyAsBase64(b64.InitialKey, invalid)
			return err
		}},
		{name: "invalid current key for encrypt pin", call: func() error {
			_, err := EncryptPinAsBase64(invalid, pin, pan, formatVersion)
			return err
		}},
		{name: "invalid current key for decrypt pin", call: func() error {
			_, err := DecryptPinAsBase64(invalid, b64.PinEnc, pan, formatVersion)
			return err
		}},
		{name: "invalid pin ciphertext", call: func() error {
			_, err := DecryptPinAsBase64(b64.CurrentKey, invalid, pan, formatVersion)
			return err
		}},
		{name: "invalid current key for encrypt data", call: func() error {
			_, err := EncryptDataAsBase64(invalid, "", data, pkg.ActionRequest)
			return err
		}},
		{name: "invalid iv character for encrypt data", call: func() error {
			_, err := EncryptDataAsBase64(b64.CurrentKey, invalid, data, pkg.ActionRequest)
			return err
		}},
		{name: "invalid current key for decrypt data", call: func() error {
			_, err := DecryptDataAsBase64(invalid, b64.DataReqEnc, "", pkg.ActionRequest)
			return err
		}},
		{name: "invalid data ciphertext", call: func() error {
			_, err := DecryptDataAsBase64(b64.CurrentKey, invalid, "", pkg.ActionRequest)
			return err
		}},
		{name: "invalid iv character for decrypt data", call: func() error {
			_, err := DecryptDataAsBase64(b64.CurrentKey, b64.DataReqEnc, invalid, pkg.ActionRequest)
			return err
		}},
		{name: "invalid current key for generate mac", call: func() error {
			_, err := GenerateMacAsBase64(invalid, data, pkg.ActionRequest)
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

// Errors moov-io raises on well formed base64 have to reach the caller too.
func TestBase64MoovErrors(t *testing.T) {
	t.Parallel()

	b64 := newSequenceBase64Item(moovInitialSequence[0])
	shortBdk := base64.StdEncoding.EncodeToString(moovInitialSequence[0].Bdk[:len(moovInitialSequence[0].Bdk)-1])

	tests := []struct {
		name string
		call func() error
	}{
		{name: "short bdk", call: func() error {
			_, err := DerivationOfInitialKeyAsBase64(shortBdk, b64.Ksn)
			return err
		}},
		{name: "unknown pin format for encrypt", call: func() error {
			_, err := EncryptPinAsBase64(b64.CurrentKey, pin, pan, "ISO-9")
			return err
		}},
		{name: "unknown pin format for decrypt", call: func() error {
			_, err := DecryptPinAsBase64(b64.CurrentKey, b64.PinEnc, pan, "ISO-9")
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

// The base64 form of [TestIVLength]: an empty IV is the zero vector and anything
// that is not one block is rejected before it reaches moov-io.
func TestBase64IVLength(t *testing.T) {
	t.Parallel()

	b64 := newSequenceBase64Item(moovInitialSequence[0])

	// The IVs are spelled as hex here so they line up with the raw and hex suites.
	encode := func(hexIV string) string {
		return base64.StdEncoding.EncodeToString(pkg.HexDecode(hexIV))
	}

	tests := []struct {
		name string
		iv   string
		// block is the one-block IV that iv has to behave like, empty when iv is rejected.
		block string
		// wantLen is the decoded length named in the error, for the rejected IVs.
		wantLen int
	}{
		{name: "empty", iv: "", block: encode("0000000000000000")},
		{name: "zero", iv: encode("0000000000000000"), block: encode("0000000000000000")},
		{name: "exact", iv: encode("0102030405060708"), block: encode("0102030405060708")},
		{name: "short", iv: encode("0102"), wantLen: 2},
		{name: "long", iv: encode("0102030405060708090A"), wantLen: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.block == "" {
				wantErr := fmt.Sprintf("iv must be %d bytes, got %d", desBlockLen, tt.wantLen)

				_, err := EncryptDataAsBase64(b64.CurrentKey, tt.iv, data, pkg.ActionRequest)
				require.ErrorContains(t, err, wantErr)

				_, err = DecryptDataAsBase64(b64.CurrentKey, b64.DataReqEnc, tt.iv, pkg.ActionRequest)
				require.ErrorContains(t, err, wantErr)
				return
			}

			ciphertext, err := EncryptDataAsBase64(b64.CurrentKey, tt.iv, data, pkg.ActionRequest)
			require.NoError(t, err)
			withBlock, err := EncryptDataAsBase64(b64.CurrentKey, tt.block, data, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, withBlock, ciphertext)

			plaintext, err := DecryptDataAsBase64(b64.CurrentKey, ciphertext, tt.iv, pkg.ActionRequest)
			require.NoError(t, err)
			rawPlaintext, err := base64.StdEncoding.DecodeString(plaintext)
			require.NoError(t, err)
			require.Equal(t, data, string(rawPlaintext[:len(data)]))
		})
	}
}
