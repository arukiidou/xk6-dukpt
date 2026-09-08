// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

// Package dukpt is moov-io wrapper for DUKPT (Derived Unique Key Per Transaction) cryptography operations.
package dukpt

import "github.com/moov-io/dukpt/pkg/des"

// [des.DeriveCurrentTransactionKey] port from moov-io
//
// Params:
//   - ik[ArrayBuffer] is 16 bytes initial key
//   - ksn[ArrayBuffer] is 10 bytes key serial number
//
// Return Params:
//   - result - 16 bytes transaction key
//   - err
func DeriveCurrentTransactionKey(ik, ksn []byte) ([]byte, error) {
	return des.DeriveCurrentTransactionKey(ik, ksn)
}

// [des.DerivationOfInitialKey] port from moov-io
//
// Params:
//   - ksn[ArrayBuffer] is 10 bytes key serial number
//   - bdk[ArrayBuffer] is 16 bytes base derivative Key
//
// Return Params:
//   - result - 16 bytes initial key
//   - err
func DerivationOfInitialKey(bdk, ksn []byte) ([]byte, error) {
	return des.DerivationOfInitialKey(bdk, ksn)
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
//   - result - cipher text
//   - err
func EncryptPin(currentKey []byte, pin, pan, format string) ([]byte, error) {
	return des.EncryptPin(currentKey, pin, pan, format)
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
//   - iv[ArrayBuffer] is initial vector, null for the default zero vector
//   - plainText is transaction request data
//   - action is "request" or "response"
//
// Return Params:
//   - result - encrypted data, zero padded to a multiple of 8 bytes
//   - err
func EncryptData(currentKey, iv []byte, plainText, action string) ([]byte, error) {
	return des.EncryptData(currentKey, iv, plainText, action)
}

// [des.DecryptData] port from moov-io
//
// Params:
//   - currentKey[ArrayBuffer] is 16 bytes transaction key
//   - ciphertext[ArrayBuffer] is encrypted text
//   - iv[ArrayBuffer] is initial vector, null for the default zero vector
//   - action is "request" or "response"
//
// Return Params:
//   - result - transaction request data, zero padded to a multiple of 8 bytes
//   - err
func DecryptData(currentKey, ciphertext, iv []byte, action string) (string, error) {
	return des.DecryptData(currentKey, ciphertext, normalizeIV(iv), action)
}

// [des.GenerateMac] port from moov-io
//
// Params:
//   - currentKey[ArrayBuffer] is 16 bytes transaction key
//   - plainText is transaction request data
//   - action is "request" or "response"
//
// Return Params:
//   - result - 8 bytes mac, use the first 4 bytes as the ANSI X9.24-1 mac
//   - err
func GenerateMac(currentKey []byte, plainText, action string) ([]byte, error) {
	return des.GenerateMac(currentKey, plainText, action)
}

// des.DecryptData widens the IV only when it is nil, so a zero-length one
// reaches cipher.NewCBCDecrypter and panics.
func normalizeIV(iv []byte) []byte {
	if len(iv) == 0 {
		return nil
	}
	return iv
}
