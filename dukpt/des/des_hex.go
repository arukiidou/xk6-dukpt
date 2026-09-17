// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package des

import (
	"encoding/hex"
	"strings"

	"github.com/moov-io/dukpt/pkg/des"
)

// encodeUpperHex returns uppercase hex, the notation used by DUKPT test vectors.
func encodeUpperHex(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}

// [des.DeriveCurrentTransactionKey] port from moov-io, as a Hex variant
//
// ik hex string - 16 bytes initial key.
// ksn hex string - 10 bytes key serial number.
//
// Return Params:
//   - result is uppercase hex string - 16 bytes transaction key
//   - err
func DeriveCurrentTransactionKeyAsHex(ik, ksn string) (string, error) {
	rawIk, err := hex.DecodeString(ik)
	if err != nil {
		return "", err
	}
	rawKsn, err := hex.DecodeString(ksn)
	if err != nil {
		return "", err
	}

	ck, err := des.DeriveCurrentTransactionKey(rawIk, rawKsn)
	if err != nil {
		return "", err
	}
	return encodeUpperHex(ck), nil
}

// [des.DerivationOfInitialKey] port from moov-io, as a Hex variant
//
// bdk hex string - 16 bytes base derivation key.
// ksn hex string - 10 bytes key serial number.
//
// Return Params:
//   - result is uppercase hex string - 16 bytes initial key
//   - err
func DerivationOfInitialKeyAsHex(bdk, ksn string) (string, error) {
	rawBdk, err := hex.DecodeString(bdk)
	if err != nil {
		return "", err
	}
	rawKsn, err := hex.DecodeString(ksn)
	if err != nil {
		return "", err
	}

	rawIk, err := des.DerivationOfInitialKey(rawBdk, rawKsn)
	if err != nil {
		return "", err
	}
	return encodeUpperHex(rawIk), nil
}

// [des.EncryptPin] port from moov-io, as a Hex variant
//
// currentKey hex string - 16 bytes transaction key.
// pin is not formatted pin string.
// pan is not formatted pan string.
// format is pinblock format
// ("ISO-0", "ISO-1", "ISO-2", "ISO-3", "ISO-4", "ANSI", "ECI1", "ECI2", "ECI3", "ECI4", "VISA1", "VISA2", "VISA3", "VISA4").
//
// Return Params:
//   - result is uppercase hex string - cipher text
//   - err
func EncryptPinAsHex(currentKey, pin, pan, format string) (string, error) {
	rawCurrentKey, err := hex.DecodeString(currentKey)
	if err != nil {
		return "", err
	}

	ciphertext, err := des.EncryptPin(rawCurrentKey, pin, pan, format)
	if err != nil {
		return "", err
	}
	return encodeUpperHex(ciphertext), nil
}

// [des.DecryptPin] port from moov-io, as a Hex variant
//
// currentKey hex string - 16 bytes transaction key.
// ciphertext hex string - encrypted pin block.
// pan is not formatted pan string.
// format is pinblock format
// ("ISO-0", "ISO-1", "ISO-2", "ISO-3", "ISO-4", "ANSI", "ECI1", "ECI2", "ECI3", "ECI4", "VISA1", "VISA2", "VISA3", "VISA4").
//
// Return Params:
//   - result is pin string (plain text)
//   - err
func DecryptPinAsHex(currentKey, ciphertext, pan, format string) (string, error) {
	rawCurrentKey, err := hex.DecodeString(currentKey)
	if err != nil {
		return "", err
	}
	rawCiphertext, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	return des.DecryptPin(rawCurrentKey, rawCiphertext, pan, format)
}

// [des.EncryptData] port from moov-io, as a Hex variant
//
// currentKey hex string - 16 bytes transaction key.
// iv hex string - initial vector, empty for the default zero vector.
// plainText is transaction request data.
// action is "request" or "response".
//
// Return Params:
//   - result is uppercase hex string - encrypted data, zero padded to a multiple of 8 bytes
//   - err
func EncryptDataAsHex(currentKey, iv, plainText, action string) (string, error) {
	rawCurrentKey, err := hex.DecodeString(currentKey)
	if err != nil {
		return "", err
	}
	rawIv, err := hex.DecodeString(iv)
	if err != nil {
		return "", err
	}
	nIv, err := normalizeIV(rawIv)
	if err != nil {
		return "", err
	}

	ciphertext, err := des.EncryptData(rawCurrentKey, nIv, plainText, action)
	if err != nil {
		return "", err
	}
	return encodeUpperHex(ciphertext), nil
}

// [des.DecryptData] port from moov-io, as a Hex variant
//
// currentKey hex string - 16 bytes transaction key.
// ciphertext hex string - encrypted text.
// iv hex string - initial vector, empty for the default zero vector.
// action is "request" or "response".
//
// Return Params:
//   - result is uppercase hex string - transaction request data, zero padded to a multiple of 8 bytes
//   - err
func DecryptDataAsHex(currentKey, ciphertext, iv, action string) (string, error) {
	rawCurrentKey, err := hex.DecodeString(currentKey)
	if err != nil {
		return "", err
	}
	rawCiphertext, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	rawIv, err := hex.DecodeString(iv)
	if err != nil {
		return "", err
	}
	nIv, err := normalizeIV(rawIv)
	if err != nil {
		return "", err
	}

	plaintext, err := des.DecryptData(rawCurrentKey, rawCiphertext, nIv, action)
	if err != nil {
		return "", err
	}
	return encodeUpperHex([]byte(plaintext)), nil
}

// [des.GenerateMac] port from moov-io, as a Hex variant
//
// currentKey hex string - 16 bytes transaction key.
// plainText is transaction request data.
// action is "request" or "response".
//
// Return Params:
//   - result is uppercase hex string - 8 bytes mac, use the first 4 bytes (8 hex characters) as the ANSI X9.24-1 mac
//   - err
func GenerateMacAsHex(currentKey, plainText, action string) (string, error) {
	rawCurrentKey, err := hex.DecodeString(currentKey)
	if err != nil {
		return "", err
	}

	mac, err := des.GenerateMac(rawCurrentKey, plainText, action)
	if err != nil {
		return "", err
	}
	return encodeUpperHex(mac), nil
}
