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
// ksn hex string - 10 bytes key serial number, its transaction counter
// must be non-zero with at most 10 one-bits. Stricter than upstream (moov-io).
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

	if err := validateKey("ik", rawIk); err != nil {
		return "", err
	}
	if err := validateKsn(rawKsn); err != nil {
		return "", err
	}
	if err := validateTransactionCounter(rawKsn); err != nil {
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

	if err := validateKey("bdk", rawBdk); err != nil {
		return "", err
	}
	if err := validateKsn(rawKsn); err != nil {
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
// cleartextPin is not formatted pin string.
// pan is not formatted pan string.
// format is pinblock format
// ("ISO-0", "ISO-1", "ISO-2", "ISO-3", "ISO-4", "ANSI", "ECI1", "ECI2", "ECI3", "ECI4", "VISA1", "VISA2", "VISA3", "VISA4").
//
// Note: This extension does not validate specific PIN block formats,
// caller should control which format is used.
// For example: ISO-4 (AES based) is obviously not usable with DES.
//
// Return Params:
//   - result is uppercase hex string - cipher text
//   - err
func EncryptPinAsHex(currentKey, cleartextPin, pan, format string) (string, error) {
	rawCurrentKey, err := hex.DecodeString(currentKey)
	if err != nil {
		return "", err
	}

	if err := validateKey("currentKey", rawCurrentKey); err != nil {
		return "", err
	}

	ciphertextPin, err := des.EncryptPin(rawCurrentKey, cleartextPin, pan, format)
	if err != nil {
		return "", err
	}
	return encodeUpperHex(ciphertextPin), nil
}

// [des.DecryptPin] port from moov-io, as a Hex variant
//
// currentKey hex string - 16 bytes transaction key.
// ciphertextPin hex string - 8 bytes encrypted pin block.
// pan is not formatted pan string.
// format is pinblock format
// ("ISO-0", "ISO-1", "ISO-2", "ISO-3", "ISO-4", "ANSI", "ECI1", "ECI2", "ECI3", "ECI4", "VISA1", "VISA2", "VISA3", "VISA4").
//
// Note: This extension does not validate specific PIN block formats,
// caller should control which format is used.
// For example: ISO-4 (AES based) is obviously not usable with DES.
//
// Return Params:
//   - result is pin string (plain text)
//   - err
func DecryptPinAsHex(currentKey, ciphertextPin, pan, format string) (string, error) {
	rawCurrentKey, err := hex.DecodeString(currentKey)
	if err != nil {
		return "", err
	}
	rawCiphertextPin, err := hex.DecodeString(ciphertextPin)
	if err != nil {
		return "", err
	}

	if err := validateKey("currentKey", rawCurrentKey); err != nil {
		return "", err
	}
	if err := validatePinBlock(rawCiphertextPin); err != nil {
		return "", err
	}

	return des.DecryptPin(rawCurrentKey, rawCiphertextPin, pan, format)
}

// [des.EncryptData] port from moov-io, as a Hex variant
//
// currentKey hex string - 16 bytes transaction key.
// iv hex string - initial vector, empty for the default zero vector.
// plainText hex string - transaction request data, must not be empty.
// action is "request" or "response", anything else is an error.
// Stricter than upstream (moov-io).
//
// Return Params:
//   - result is uppercase hex string - encrypted data, zero padded to a multiple of 8 bytes
//   - err
func EncryptDataAsHex(currentKey, iv, plainText, action string) (string, error) {
	if err := validateAction(action); err != nil {
		return "", err
	}
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
	rawPlainText, err := hex.DecodeString(plainText)
	if err != nil {
		return "", err
	}

	if err := validateKey("currentKey", rawCurrentKey); err != nil {
		return "", err
	}

	if err := validatePlainText(rawPlainText); err != nil {
		return "", err
	}

	ciphertext, err := des.EncryptData(rawCurrentKey, nIv, string(rawPlainText), action)
	if err != nil {
		return "", err
	}
	return encodeUpperHex(ciphertext), nil
}

// [des.DecryptData] port from moov-io, as a Hex variant
//
// currentKey hex string - 16 bytes transaction key.
// ciphertext hex string - encrypted text, a non-empty multiple of 8 bytes.
// Stricter than upstream (moov-io), which zero pads other lengths.
// iv hex string - initial vector, empty for the default zero vector.
// action is "request" or "response", anything else is an error.
// Stricter than upstream (moov-io).
//
// Return Params:
//   - result is uppercase hex string - transaction request data, zero padded to a multiple of 8 bytes
//   - err
func DecryptDataAsHex(currentKey, ciphertext, iv, action string) (string, error) {
	if err := validateAction(action); err != nil {
		return "", err
	}
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

	if err := validateKey("currentKey", rawCurrentKey); err != nil {
		return "", err
	}

	if err := validateCiphertext(rawCiphertext); err != nil {
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
// plainText hex string - transaction request data, must not be empty.
// action is "request" or "response", anything else is an error.
// Stricter than upstream (moov-io).
//
// Return Params:
//   - result is uppercase hex string - 8 bytes mac, use the first 4 bytes (8 hex characters) as the ANSI X9.24-1 mac
//   - err
func GenerateMacAsHex(currentKey, plainText, action string) (string, error) {
	if err := validateAction(action); err != nil {
		return "", err
	}
	rawCurrentKey, err := hex.DecodeString(currentKey)
	if err != nil {
		return "", err
	}
	rawPlainText, err := hex.DecodeString(plainText)
	if err != nil {
		return "", err
	}

	if err := validateKey("currentKey", rawCurrentKey); err != nil {
		return "", err
	}

	if err := validatePlainText(rawPlainText); err != nil {
		return "", err
	}

	mac, err := des.GenerateMac(rawCurrentKey, string(rawPlainText), action)
	if err != nil {
		return "", err
	}
	return encodeUpperHex(mac), nil
}
