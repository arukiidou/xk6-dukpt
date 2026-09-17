// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package des

import (
	"fmt"
	"testing"

	"github.com/grafana/sobek"
	"github.com/moov-io/dukpt/pkg"
	"github.com/moov-io/dukpt/pkg/des"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/v2/js/modulestest"
)

// newTestModule returns a module backed by a test runtime, plus an unwrap
// helper that flattens its ArrayBuffer results into raw bytes.
func newTestModule(t *testing.T) (*module, func(*sobek.ArrayBuffer, error) []byte) {
	t.Helper()

	instance, ok := new(DesModule).NewModuleInstance(modulestest.NewRuntime(t).VU).(*module)
	require.True(t, ok)

	unwrap := func(ab *sobek.ArrayBuffer, err error) []byte {
		t.Helper()

		require.NoError(t, err)
		require.NotNil(t, ab)
		return ab.Bytes()
	}
	return instance, unwrap
}

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
var moovInitialSequence = []SequenceItem{
	{
		Bdk:         moovTestBdk,
		Ksn:         pkg.HexDecode("FFFF9876543210E00001"),                             // //+YdlQyEOAAAQ==
		InitialKey:  pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"),                 // asKS+qExW02FirOj19WTOg==
		CurrentKey:  pkg.HexDecode("042666B49184CFA368DE9628D0397BC9"),                 // BCZmtJGEz6No3pYo0Dl7yQ==
		PinEnc:      pkg.HexDecode("1B9C1845EB993A7A"),                                 // G5wYReuZOno=
		DataReqEnc:  pkg.HexDecode("FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F"), // /A1Tt+of2p7miq8ucNm5UGIpviqpk/BP
		DataResEnc:  pkg.HexDecode("1FCC89AF66222F27B903898BB2BC8589CDBFDE5EC6AFCC25"), // H8yJr2YiLye5A4mLsryFic2/3l7Gr8wl
		RequestMac:  pkg.HexDecode("9CCC78173FC4FB64"),                                 // nMx4Fz/E+2Q=
		ResponseMac: pkg.HexDecode("20364223C1FF00FA"),                                 // IDZCI8H/APo=
	},
	{
		Bdk:         moovTestBdk,
		Ksn:         pkg.HexDecode("FFFF9876543210E00002"),                             // //+YdlQyEOAAAg==
		InitialKey:  pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"),                 // asKS+qExW02FirOj19WTOg==
		CurrentKey:  pkg.HexDecode("C46551CEF9FD24B0AA9AD834130D3BC7"),                 // xGVRzvn9JLCqmtg0Ew07xw==
		PinEnc:      pkg.HexDecode("10A01C8D02C69107"),                                 // EKAcjQLGkQc=
		DataReqEnc:  pkg.HexDecode("A2B4E70F846E63D68775B7215EB4563DFD3037244C61CC13"), // orTnD4RuY9aHdbchXrRWPf0wNyRMYcwT
		DataResEnc:  pkg.HexDecode("5B692A6B1FDD5E25B0DEFAFDE1672E402F8011360CFF3508"), // W2kqax/dXiWw3vr94WcuQC+AETYM/zUI
		RequestMac:  pkg.HexDecode("F608A9BCA6FFC311"),                                 // 9gipvKb/wxE=
		ResponseMac: pkg.HexDecode("D1FCA6BEF05D24D2"),                                 // 0fymvvBdJNI=
	},
	{
		Bdk:         moovTestBdk,
		Ksn:         pkg.HexDecode("FFFF9876543210E00003"),                             // //+YdlQyEOAAAw==
		InitialKey:  pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"),                 // asKS+qExW02FirOj19WTOg==
		CurrentKey:  pkg.HexDecode("0DF3D9422ACA56E547676D07AD6BADFA"),                 // DfPZQirKVuVHZ20HrWut+g==
		PinEnc:      pkg.HexDecode("18DC07B94797B466"),                                 // GNwHuUeXtGY=
		DataReqEnc:  pkg.HexDecode("BD751E65F10E75B6C1D5B1D283496A36C2DE21D993C387A7"), // vXUeZfEOdbbB1bHSg0lqNsLeIdmTw4en
		DataResEnc:  pkg.HexDecode("345992D4163E4926C927BFD8ABD5D76F087A9CE81D5A27B8"), // NFmS1BY+SSbJJ7/Yq9XXbwh6nOgdWie4
		RequestMac:  pkg.HexDecode("20B59A4FEAC937E3"),                                 // ILWaT+rJN+M=
		ResponseMac: pkg.HexDecode("BAD4CC9CC2AE326C"),                                 // utTMnMKuMmw=
	},
	{
		Bdk:         moovTestBdk,
		Ksn:         pkg.HexDecode("FFFF9876543210E00004"),                             // //+YdlQyEOAABA==
		InitialKey:  pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"),                 // asKS+qExW02FirOj19WTOg==
		CurrentKey:  pkg.HexDecode("279C0F6AEED0BE652B2C733E1383AE91"),                 // J5wPau7QvmUrLHM+E4OukQ==
		PinEnc:      pkg.HexDecode("0BC79509D5645DF7"),                                 // C8eVCdVkXfc=
		DataReqEnc:  pkg.HexDecode("1118F50947441BBDA3C8C70220021A12EC31CC473F7215F4"), // ERj1CUdEG72jyMcCIAIaEuwxzEc/chX0
		DataResEnc:  pkg.HexDecode("418C7413576C0D1819E785D3807AF32334231FDEC23414DB"), // QYx0E1dsDRgZ54XTgHrzIzQjH97CNBTb
		RequestMac:  pkg.HexDecode("C7BFA6CC44161828"),                                 // x7+mzEQWGCg=
		ResponseMac: pkg.HexDecode("1EB08AEECE6FF0C2"),                                 // HrCK7s5v8MI=
	},
	{
		Bdk:         moovTestBdk,
		Ksn:         pkg.HexDecode("FFFF9876543210E00005"),                             // //+YdlQyEOAABQ==
		InitialKey:  pkg.HexDecode("6AC292FAA1315B4D858AB3A3D7D5933A"),                 // asKS+qExW02FirOj19WTOg==
		CurrentKey:  pkg.HexDecode("5F8DC6D2C845C125508DDC048093B83F"),                 // X43G0shFwSVQjdwEgJO4Pw==
		PinEnc:      pkg.HexDecode("5BC0AF22AD87B327"),                                 // W8CvIq2Hsyc=
		DataReqEnc:  pkg.HexDecode("9FD7BD1EC28845ACA93367A9DA9317BD555C6B33AE22D365"), // n9e9HsKIRaypM2ep2pMXvVVcazOuItNl
		DataResEnc:  pkg.HexDecode("7D4C109E49E83355A556AE949EED359F4404E7A2F0167C00"), // fUwQnknoM1WlVq6Unu01n0QE56LwFnwA
		RequestMac:  pkg.HexDecode("0202B96339022058"),                                 // AgK5YzkCIFg=
		ResponseMac: pkg.HexDecode("5CBE3E81D1D2A0FB"),                                 // XL4+gdHSoPs=
	},
}

func TestMoovCompatibility(t *testing.T) {
	t.Parallel()

	for index, moov := range moovInitialSequence {
		t.Run(fmt.Sprintf("Sequence #%d KSN: %s", index+1, pkg.HexEncode(moov.Ksn)), func(t *testing.T) {
			t.Parallel()

			m, unwrap := newTestModule(t)

			ik := unwrap(m.DerivationOfInitialKey(moov.Bdk, moov.Ksn))
			desIk, err := des.DerivationOfInitialKey(moov.Bdk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, desIk, ik)
			require.Equal(t, moov.InitialKey, ik)

			ck := unwrap(m.DeriveCurrentTransactionKey(ik, moov.Ksn))
			desCk, err := des.DeriveCurrentTransactionKey(desIk, moov.Ksn)
			require.NoError(t, err)
			require.Equal(t, desCk, ck)
			require.Equal(t, moov.CurrentKey, ck)

			pinEnc := unwrap(m.EncryptPin(ck, pin, pan, formatVersion))
			desPinEnc, err := des.EncryptPin(desCk, pin, pan, formatVersion)
			require.NoError(t, err)
			require.Equal(t, desPinEnc, pinEnc)
			require.Equal(t, moov.PinEnc, pinEnc)

			decPin, err := DecryptPin(ck, pinEnc, pan, formatVersion)
			require.NoError(t, err)
			require.Equal(t, pin, decPin)

			reqEnc := unwrap(m.EncryptData(ck, nil, data, pkg.ActionRequest))
			require.Equal(t, moov.DataReqEnc, reqEnc)

			decReq := unwrap(m.DecryptData(ck, reqEnc, nil, pkg.ActionRequest))
			require.Equal(t, data, string(decReq[:len(data)]))

			resEnc := unwrap(m.EncryptData(ck, nil, data, pkg.ActionResponse))
			require.Equal(t, moov.DataResEnc, resEnc)

			decRes := unwrap(m.DecryptData(ck, resEnc, nil, pkg.ActionResponse))
			require.Equal(t, data, string(decRes[:len(data)]))

			// An empty IV must behave like the nil default instead of panicking.
			decEmptyIV := unwrap(m.DecryptData(ck, reqEnc, []byte{}, pkg.ActionRequest))
			require.Equal(t, decReq, decEmptyIV)

			reqMac := unwrap(m.GenerateMac(ck, data, pkg.ActionRequest))
			desReqMac, err := des.GenerateMac(desCk, data, pkg.ActionRequest)
			require.NoError(t, err)
			require.Equal(t, desReqMac, reqMac)
			require.Equal(t, moov.RequestMac, reqMac)

			resMac := unwrap(m.GenerateMac(ck, data, pkg.ActionResponse))
			desResMac, err := des.GenerateMac(desCk, data, pkg.ActionResponse)
			require.NoError(t, err)
			require.Equal(t, desResMac, resMac)
			require.Equal(t, moov.ResponseMac, resMac)
		})
	}
}

// An IV that is not one block panics in moov-io, so both sides must reject it.
// A missing one still has to fall back to the zero vector.
func TestIVLength(t *testing.T) {
	t.Parallel()

	currentKey := pkg.HexDecode("042666B49184CFA368DE9628D0397BC9")
	zeroBlock := pkg.HexDecode("0000000000000000")

	tests := []struct {
		name string
		iv   []byte
		// block is the one-block IV that iv has to behave like, nil when iv is rejected.
		block []byte
	}{
		{name: "nil", iv: nil, block: zeroBlock},
		{name: "empty", iv: []byte{}, block: zeroBlock},
		{name: "exact", iv: pkg.HexDecode("0102030405060708"), block: pkg.HexDecode("0102030405060708")},
		{name: "short", iv: pkg.HexDecode("0102")},
		{name: "long", iv: pkg.HexDecode("0102030405060708090A")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, unwrap := newTestModule(t)

			if tt.block == nil {
				_, err := m.EncryptData(currentKey, tt.iv, data, pkg.ActionRequest)
				require.ErrorContains(t, err, fmt.Sprintf("iv must be 8 bytes, got %d", len(tt.iv)))

				_, err = m.DecryptData(currentKey, nil, tt.iv, pkg.ActionRequest)
				require.ErrorContains(t, err, fmt.Sprintf("iv must be 8 bytes, got %d", len(tt.iv)))
				return
			}

			ciphertext := unwrap(m.EncryptData(currentKey, tt.iv, data, pkg.ActionRequest))
			require.Equal(t, unwrap(m.EncryptData(currentKey, tt.block, data, pkg.ActionRequest)), ciphertext)

			plaintext := unwrap(m.DecryptData(currentKey, ciphertext, tt.iv, pkg.ActionRequest))
			require.Equal(t, data, string(plaintext[:len(data)]))
		})
	}
}

// Errors moov-io raises have to reach the caller instead of an empty buffer.
//
// Only a short BDK and an unknown PIN block format are reachable: moov-io copies
// every other key into a fixed 16 byte buffer, zero padding it, so a wrong length
// ik or currentKey silently produces a wrong result rather than an error. That is
// why the error returns of DeriveCurrentTransactionKey and GenerateMac stay
// uncovered.
func TestRawMoovErrors(t *testing.T) {
	t.Parallel()

	seq := moovInitialSequence[0]

	tests := []struct {
		name string
		call func(*module) error
	}{
		{name: "short bdk", call: func(m *module) error {
			_, err := m.DerivationOfInitialKey(seq.Bdk[:len(seq.Bdk)-1], seq.Ksn)
			return err
		}},
		{name: "unknown pin format for encrypt", call: func(m *module) error {
			_, err := m.EncryptPin(seq.CurrentKey, pin, pan, "ISO-9")
			return err
		}},
		{name: "unknown pin format for decrypt", call: func(_ *module) error {
			_, err := DecryptPin(seq.CurrentKey, seq.PinEnc, pan, "ISO-9")
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, _ := newTestModule(t)
			require.Error(t, tt.call(m))
		})
	}
}

// The raw APIs hand their results to JS as an ArrayBuffer backed by the same
// bytes, so a moov-io call that writes into one of its arguments would be
// visible from the script. None of them may touch what they are given.
func TestRawInputsNotModified(t *testing.T) {
	t.Parallel()

	seq := moovInitialSequence[0]
	iv := pkg.HexDecode("0102030405060708")

	tests := []struct {
		name string
		// inputs are the arguments that have to come back unchanged.
		inputs [][]byte
		call   func(*module, func(*sobek.ArrayBuffer, error) []byte)
	}{
		{name: "DerivationOfInitialKey", inputs: [][]byte{seq.Bdk, seq.Ksn}, call: func(m *module, unwrap func(*sobek.ArrayBuffer, error) []byte) {
			unwrap(m.DerivationOfInitialKey(seq.Bdk, seq.Ksn))
		}},
		{name: "DeriveCurrentTransactionKey", inputs: [][]byte{seq.InitialKey, seq.Ksn}, call: func(m *module, unwrap func(*sobek.ArrayBuffer, error) []byte) {
			unwrap(m.DeriveCurrentTransactionKey(seq.InitialKey, seq.Ksn))
		}},
		{name: "EncryptPin", inputs: [][]byte{seq.CurrentKey}, call: func(m *module, unwrap func(*sobek.ArrayBuffer, error) []byte) {
			unwrap(m.EncryptPin(seq.CurrentKey, pin, pan, formatVersion))
		}},
		{name: "DecryptPin", inputs: [][]byte{seq.CurrentKey, seq.PinEnc}, call: func(m *module, _ func(*sobek.ArrayBuffer, error) []byte) {
			_, err := DecryptPin(seq.CurrentKey, seq.PinEnc, pan, formatVersion)
			require.NoError(t, err)
		}},
		{name: "EncryptData", inputs: [][]byte{seq.CurrentKey, iv}, call: func(m *module, unwrap func(*sobek.ArrayBuffer, error) []byte) {
			unwrap(m.EncryptData(seq.CurrentKey, iv, data, pkg.ActionRequest))
		}},
		{name: "DecryptData", inputs: [][]byte{seq.CurrentKey, seq.DataReqEnc, iv}, call: func(m *module, unwrap func(*sobek.ArrayBuffer, error) []byte) {
			unwrap(m.DecryptData(seq.CurrentKey, seq.DataReqEnc, iv, pkg.ActionRequest))
		}},
		{name: "GenerateMac", inputs: [][]byte{seq.CurrentKey}, call: func(m *module, unwrap func(*sobek.ArrayBuffer, error) []byte) {
			unwrap(m.GenerateMac(seq.CurrentKey, data, pkg.ActionRequest))
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, unwrap := newTestModule(t)

			// The fixtures are shared, so compare against a copy taken up front.
			want := make([][]byte, len(tt.inputs))
			for i, in := range tt.inputs {
				want[i] = append([]byte(nil), in...)
			}

			tt.call(m, unwrap)

			for i, in := range tt.inputs {
				require.Equal(t, want[i], in, "input must not be modified")
			}
		})
	}
}
