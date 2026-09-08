// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"encoding/base64"

	"github.com/moov-io/dukpt/pkg/des"
)

// [des.DeriveCurrentTransactionKey] port from moov-io
//
// ik base64 string - 16 bytes initial key.
// ksn base64 string - 10 bytes key serial number.
//
// Return Params:
//   - result is base64 string - 16 bytes transaction key
//   - err
func DeriveCurrentTransactionKeyAsBase64(ik, ksn string) (string, error) {
	rawIk, err := base64.StdEncoding.DecodeString(ik)
	if err != nil {
		return "", err
	}
	rawKsn, err := base64.StdEncoding.DecodeString(ksn)
	if err != nil {
		return "", err
	}

	ck, err := des.DeriveCurrentTransactionKey(rawIk, rawKsn)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ck), nil
}

// [des.DerivationOfInitialKey] port from moov-io
//
// bdk base64 string - 16 bytes base derivation key.
// ksn base64 string - 10 bytes key serial number.
//
// Return Params:
//   - result is base64 string - 16 bytes initial key
//   - err
func DerivationOfInitialKeyAsBase64(bdk, ksn string) (string, error) {
	rawBdk, err := base64.StdEncoding.DecodeString(bdk)
	if err != nil {
		return "", err
	}
	rawKsn, err := base64.StdEncoding.DecodeString(ksn)
	if err != nil {
		return "", err
	}

	rawCk, err := des.DerivationOfInitialKey(rawBdk, rawKsn)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(rawCk), nil
}

// [des.EncryptPin] port from moov-io
//
// currentKey base64 string - 16 bytes transaction key.
// pin is not formatted pin string.
// pan is not formatted pan string.
// format is pinblock format
// ("ISO-0", "ISO-1", "ISO-2", "ISO-3", "ISO-4", "ANSI", "ECI1", "ECI2", "ECI3", "ECI4", "VISA1", "VISA2", "VISA3", "VISA4").
//
// Return Params:
//   - result is base64 string - cipher text
//   - err
func EncryptPinAsBase64(currentKey, pin, pan, format string) (string, error) {
	rawCurrentKey, err := base64.StdEncoding.DecodeString(currentKey)
	if err != nil {
		return "", err
	}

	ciphertext, err := des.EncryptPin(rawCurrentKey, pin, pan, format)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// [des.DecryptPin] port from moov-io
//
// currentKey base64 string - 16 bytes transaction key.
// ciphertext base64 string - encrypted pin block.
// pan is not formatted pan string.
// format is pinblock format
// ("ISO-0", "ISO-1", "ISO-2", "ISO-3", "ISO-4", "ANSI", "ECI1", "ECI2", "ECI3", "ECI4", "VISA1", "VISA2", "VISA3", "VISA4").
//
// Return Params:
//   - result is pin string (plain text)
//   - err
func DecryptPinAsBase64(currentKey, ciphertext, pan, format string) (string, error) {
	rawCurrentKey, err := base64.StdEncoding.DecodeString(currentKey)
	if err != nil {
		return "", err
	}
	rawCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	return des.DecryptPin(rawCurrentKey, rawCiphertext, pan, format)
}

// [des.EncryptData] port from moov-io
//
// currentKey base64 string - 16 bytes transaction key.
// iv base64 string - initial vector, empty for the default zero vector.
// plainText is transaction request data.
// action is "request" or "response".
//
// Return Params:
//   - result is base64 string - encrypted data, zero padded to a multiple of 8 bytes
//   - err
func EncryptDataAsBase64(currentKey, iv, plainText, action string) (string, error) {
	rawCurrentKey, err := base64.StdEncoding.DecodeString(currentKey)
	if err != nil {
		return "", err
	}
	rawIv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", err
	}

	ciphertext, err := des.EncryptData(rawCurrentKey, normalizeIV(rawIv), plainText, action)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// [des.DecryptData] port from moov-io
//
// currentKey base64 string - 16 bytes transaction key.
// ciphertext base64 string - encrypted text.
// iv base64 string - initial vector, empty for the default zero vector.
// action is "request" or "response".
//
// Return Params:
//   - result is transaction request data (plain text), zero padded to a multiple of 8 bytes
//   - err
func DecryptDataAsBase64(currentKey, ciphertext, iv, action string) (string, error) {
	rawCurrentKey, err := base64.StdEncoding.DecodeString(currentKey)
	if err != nil {
		return "", err
	}
	rawCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	rawIv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", err
	}

	return des.DecryptData(rawCurrentKey, rawCiphertext, normalizeIV(rawIv), action)
}

// [des.GenerateMac] port from moov-io
//
// currentKey base64 string - 16 bytes transaction key.
// plainText is transaction request data.
// action is "request" or "response".
//
// Return Params:
//   - result is base64 string - 8 bytes mac, use the first 4 bytes as the ANSI X9.24-1 mac
//   - err
func GenerateMacAsBase64(currentKey, plainText, action string) (string, error) {
	rawCurrentKey, err := base64.StdEncoding.DecodeString(currentKey)
	if err != nil {
		return "", err
	}

	mac, err := des.GenerateMac(rawCurrentKey, plainText, action)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac), nil
}
