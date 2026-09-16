// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package des

import (
	"fmt"

	"github.com/grafana/sobek"
	"github.com/moov-io/dukpt/pkg/des"
)

// desBlockLen is the TDEA block length, the only IV length CBC mode accepts.
const desBlockLen = 8

// newArrayBuffer returns raw bytes as JS ArrayBuffer.
func (m *module) newArrayBuffer(b []byte) *sobek.ArrayBuffer {
	ab := m.vu.Runtime().NewArrayBuffer(b)
	return &ab
}

// [des.DeriveCurrentTransactionKey] port from moov-io
//
// Params:
//   - ik[ArrayBuffer] is 16 bytes initial key
//   - ksn[ArrayBuffer] is 10 bytes key serial number
//
// Return Params:
//   - result[ArrayBuffer] - 16 bytes transaction key
//   - err
func (m *module) DeriveCurrentTransactionKey(ik, ksn []byte) (*sobek.ArrayBuffer, error) {
	ck, err := des.DeriveCurrentTransactionKey(ik, ksn)
	if err != nil {
		return nil, err
	}
	return m.newArrayBuffer(ck), nil
}

// [des.DerivationOfInitialKey] port from moov-io
//
// Params:
//   - ksn[ArrayBuffer] is 10 bytes key serial number
//   - bdk[ArrayBuffer] is 16 bytes base derivative Key
//
// Return Params:
//   - result[ArrayBuffer] - 16 bytes initial key
//   - err
func (m *module) DerivationOfInitialKey(bdk, ksn []byte) (*sobek.ArrayBuffer, error) {
	ik, err := des.DerivationOfInitialKey(bdk, ksn)
	if err != nil {
		return nil, err
	}
	return m.newArrayBuffer(ik), nil
}

// [des.EncryptPin] port from moov-io
//
// Params:
//   - currentKey[ArrayBuffer] is 16 bytes transaction key
//   - pin is not formatted pin string
//   - pan is not formatted pan string
//   - format is pinblock format
//     ("ISO-0", "ISO-1", "ISO-2", "ISO-3", "ISO-4", "ANSI", "ECI1", "ECI2", "ECI3", "ECI4", "VISA1", "VISA2", "VISA3", "VISA4")
//
// Return Params:
//   - result[ArrayBuffer] - cipher text
//   - err
func (m *module) EncryptPin(currentKey []byte, pin, pan, format string) (*sobek.ArrayBuffer, error) {
	ciphertext, err := des.EncryptPin(currentKey, pin, pan, format)
	if err != nil {
		return nil, err
	}
	return m.newArrayBuffer(ciphertext), nil
}

// [des.DecryptPin] port from moov-io
//
// Params:
//   - currentKey[ArrayBuffer] is 16 bytes transaction key
//   - ciphertext[ArrayBuffer] is encrypted pin block
//   - pan is not formatted pan string
//   - format is pinblock format
//     ("ISO-0", "ISO-1", "ISO-2", "ISO-3", "ISO-4", "ANSI", "ECI1", "ECI2", "ECI3", "ECI4", "VISA1", "VISA2", "VISA3", "VISA4")
//
// Return Params:
//   - result - pin string (plain text)
//   - err
func DecryptPin(currentKey, ciphertext []byte, pan, format string) (string, error) {
	return des.DecryptPin(currentKey, ciphertext, pan, format)
}

// [des.EncryptData] port from moov-io
//
// Params:
//   - currentKey[ArrayBuffer] is 16 bytes transaction key
//   - iv[ArrayBuffer] is the 8 byte initial vector, null or empty for the default zero vector
//   - plainText is transaction request data
//   - action is "request" or "response"
//
// Return Params:
//   - result[ArrayBuffer] - encrypted data, zero padded to a multiple of 8 bytes
//   - err
func (m *module) EncryptData(currentKey, iv []byte, plainText, action string) (*sobek.ArrayBuffer, error) {
	iv, err := normalizeIV(iv)
	if err != nil {
		return nil, err
	}

	ciphertext, err := des.EncryptData(currentKey, iv, plainText, action)
	if err != nil {
		return nil, err
	}
	return m.newArrayBuffer(ciphertext), nil
}

// [des.DecryptData] port from moov-io
//
// Params:
//   - currentKey[ArrayBuffer] is 16 bytes transaction key
//   - ciphertext[ArrayBuffer] is encrypted text
//   - iv[ArrayBuffer] is the 8 byte initial vector, null or empty for the default zero vector
//   - action is "request" or "response"
//
// Return Params:
//   - result - transaction request data, zero padded to a multiple of 8 bytes
//   - err
func DecryptData(currentKey, ciphertext, iv []byte, action string) (string, error) {
	iv, err := normalizeIV(iv)
	if err != nil {
		return "", err
	}
	return des.DecryptData(currentKey, ciphertext, iv, action)
}

// [des.GenerateMac] port from moov-io
//
// Params:
//   - currentKey[ArrayBuffer] is 16 bytes transaction key
//   - plainText is transaction request data
//   - action is "request" or "response"
//
// Return Params:
//   - result[ArrayBuffer] - 8 bytes mac, use the first 4 bytes as the ANSI X9.24-1 mac
//   - err
func (m *module) GenerateMac(currentKey []byte, plainText, action string) (*sobek.ArrayBuffer, error) {
	mac, err := des.GenerateMac(currentKey, plainText, action)
	if err != nil {
		return nil, err
	}
	return m.newArrayBuffer(mac), nil
}

// normalizeIV returns the IV as the one block CBC mode demands. A missing or
// empty IV becomes the zero vector; any other length is rejected, instead of
// letting moov-io pad, truncate, or panic on it.
func normalizeIV(iv []byte) ([]byte, error) {
	if len(iv) == 0 {
		return make([]byte, desBlockLen), nil
	}
	if len(iv) != desBlockLen {
		return nil, fmt.Errorf("iv must be %d bytes, got %d", desBlockLen, len(iv))
	}
	return iv, nil
}
