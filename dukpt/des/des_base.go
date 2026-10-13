// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package des

import (
	"errors"
	"fmt"
	"math/bits"

	"github.com/arukiidou/xk6-dukpt/internal/jsbytes"
	"github.com/grafana/sobek"
	"github.com/moov-io/dukpt/pkg"
	"github.com/moov-io/dukpt/pkg/des"
)

// desBlockLen is the TDEA block length, the only IV length CBC mode accepts.
const (
	desBlockLen = 8
	desKeyLen   = 16
	desKsnLen   = 10
)

// newArrayBuffer returns raw bytes as JS ArrayBuffer.
func (m *module) newArrayBuffer(b []byte) *sobek.ArrayBuffer {
	ab := m.vu.Runtime().NewArrayBuffer(b)
	return &ab
}

// [des.DeriveCurrentTransactionKey] port from moov-io
//
// Params:
//   - ik[ArrayBuffer] is 16 bytes initial key
//   - ksn[ArrayBuffer] is 10 bytes key serial number, its transaction counter
//     must be non-zero with at most 10 one-bits. Stricter than upstream (moov-io).
//
// Return Params:
//   - result[ArrayBuffer] - 16 bytes transaction key
//   - err
func (m *module) DeriveCurrentTransactionKey(ikArg, ksnArg sobek.Value) (*sobek.ArrayBuffer, error) {
	ik, err := jsbytes.Required("ik", ikArg)
	if err != nil {
		return nil, err
	}
	ksn, err := jsbytes.Required("ksn", ksnArg)
	if err != nil {
		return nil, err
	}

	if err := validateKey("ik", ik); err != nil {
		return nil, err
	}
	if err := validateKsn(ksn); err != nil {
		return nil, err
	}
	if err := validateTransactionCounter(ksn); err != nil {
		return nil, err
	}

	ck, err := des.DeriveCurrentTransactionKey(ik, ksn)
	if err != nil {
		return nil, err
	}
	return m.newArrayBuffer(ck), nil
}

// [des.DerivationOfInitialKey] port from moov-io
//
// Params:
//   - bdk[ArrayBuffer] is 16 bytes base derivation key
//   - ksn[ArrayBuffer] is 10 bytes key serial number
//
// Return Params:
//   - result[ArrayBuffer] - 16 bytes initial key
//   - err
func (m *module) DerivationOfInitialKey(bdkArg, ksnArg sobek.Value) (*sobek.ArrayBuffer, error) {
	bdk, err := jsbytes.Required("bdk", bdkArg)
	if err != nil {
		return nil, err
	}
	ksn, err := jsbytes.Required("ksn", ksnArg)
	if err != nil {
		return nil, err
	}

	if err := validateKey("bdk", bdk); err != nil {
		return nil, err
	}
	if err := validateKsn(ksn); err != nil {
		return nil, err
	}

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
//   - cleartextPin is not formatted pin string
//   - pan is not formatted pan string
//   - format is pinblock format
//     ("ISO-0", "ISO-1", "ISO-2", "ISO-3", "ISO-4", "ANSI", "ECI1", "ECI2", "ECI3", "ECI4", "VISA1", "VISA2", "VISA3", "VISA4")
//
// Note: This extension does not validate specific PIN block formats,
// caller should control which format is used.
// For example: ISO-4 (AES based) is obviously not usable with DES.
//
// Return Params:
//   - result[ArrayBuffer] - cipher text
//   - err
func (m *module) EncryptPin(currentKeyArg sobek.Value, cleartextPin, pan, format string) (*sobek.ArrayBuffer, error) {
	currentKey, err := jsbytes.Required("currentKey", currentKeyArg)
	if err != nil {
		return nil, err
	}

	if err := validateKey("currentKey", currentKey); err != nil {
		return nil, err
	}

	ciphertextPin, err := des.EncryptPin(currentKey, cleartextPin, pan, format)
	if err != nil {
		return nil, err
	}
	return m.newArrayBuffer(ciphertextPin), nil
}

// [des.DecryptPin] port from moov-io
//
// Params:
//   - currentKey[ArrayBuffer] is 16 bytes transaction key
//   - ciphertextPin[ArrayBuffer] is 8 bytes encrypted pin block
//   - pan is not formatted pan string
//   - format is pinblock format
//     ("ISO-0", "ISO-1", "ISO-2", "ISO-3", "ISO-4", "ANSI", "ECI1", "ECI2", "ECI3", "ECI4", "VISA1", "VISA2", "VISA3", "VISA4")
//
// Note: This extension does not validate specific PIN block formats,
// caller should control which format is used.
// For example: ISO-4 (AES based) is obviously not usable with DES.
//
// Return Params:
//   - result - pin string (plain text)
//   - err
func DecryptPin(currentKeyArg, ciphertextPinArg sobek.Value, pan, format string) (string, error) {
	currentKey, err := jsbytes.Required("currentKey", currentKeyArg)
	if err != nil {
		return "", err
	}
	ciphertextPin, err := jsbytes.Required("ciphertextPin", ciphertextPinArg)
	if err != nil {
		return "", err
	}

	if err := validateKey("currentKey", currentKey); err != nil {
		return "", err
	}
	if err := validatePinBlock(ciphertextPin); err != nil {
		return "", err
	}

	return des.DecryptPin(currentKey, ciphertextPin, pan, format)
}

// [des.EncryptData] port from moov-io
//
// Params:
//   - currentKey[ArrayBuffer] is 16 bytes transaction key
//   - iv[ArrayBuffer] is the 8 byte initial vector, null, undefined or empty for the default zero vector
//   - plainText[ArrayBuffer] is transaction request data, must not be empty.
//   - action is "request" or "response", anything else is an error.
//     Stricter than upstream (moov-io).
//
// Return Params:
//   - result[ArrayBuffer] - encrypted data, zero padded to a multiple of 8 bytes
//   - err
func (m *module) EncryptData(currentKeyArg, ivArg, plainTextArg sobek.Value, action string) (*sobek.ArrayBuffer, error) {
	currentKey, err := jsbytes.Required("currentKey", currentKeyArg)
	if err != nil {
		return nil, err
	}
	rawIv, err := jsbytes.Optional("iv", ivArg)
	if err != nil {
		return nil, err
	}
	rawPlainText, err := jsbytes.Required("plainText", plainTextArg)
	if err != nil {
		return nil, err
	}

	if err := validateAction(action); err != nil {
		return nil, err
	}
	iv, err := normalizeIV(rawIv)
	if err != nil {
		return nil, err
	}

	if err := validateKey("currentKey", currentKey); err != nil {
		return nil, err
	}

	if err := validatePlainText(rawPlainText); err != nil {
		return nil, err
	}

	ciphertext, err := des.EncryptData(currentKey, iv, string(rawPlainText), action)
	if err != nil {
		return nil, err
	}
	return m.newArrayBuffer(ciphertext), nil
}

// [des.DecryptData] port from moov-io
//
// Params:
//   - currentKey[ArrayBuffer] is 16 bytes transaction key
//   - ciphertext[ArrayBuffer] is encrypted text, a non-empty multiple of 8 bytes
//     Stricter than upstream (moov-io), which zero pads other lengths.
//   - iv[ArrayBuffer] is the 8 byte initial vector, null, undefined or empty for the default zero vector
//   - action is "request" or "response", anything else is an error.
//     Stricter than upstream (moov-io).
//
// Return Params:
//   - result[ArrayBuffer] - transaction request data, zero padded to a multiple of 8 bytes
//   - err
func (m *module) DecryptData(currentKeyArg, ciphertextArg, ivArg sobek.Value, action string) (*sobek.ArrayBuffer, error) {
	currentKey, err := jsbytes.Required("currentKey", currentKeyArg)
	if err != nil {
		return nil, err
	}
	ciphertext, err := jsbytes.Required("ciphertext", ciphertextArg)
	if err != nil {
		return nil, err
	}
	rawIv, err := jsbytes.Optional("iv", ivArg)
	if err != nil {
		return nil, err
	}

	if err := validateAction(action); err != nil {
		return nil, err
	}
	iv, err := normalizeIV(rawIv)
	if err != nil {
		return nil, err
	}

	if err := validateKey("currentKey", currentKey); err != nil {
		return nil, err
	}

	if err := validateCiphertext(ciphertext); err != nil {
		return nil, err
	}

	plaintext, err := des.DecryptData(currentKey, ciphertext, iv, action)
	if err != nil {
		return nil, err
	}
	// moov-io hands the plaintext back as a string, but it is the raw padded
	// block, so it goes to JS as bytes like every other result.
	return m.newArrayBuffer([]byte(plaintext)), nil
}

// [des.GenerateMac] port from moov-io
//
// Params:
//   - currentKey[ArrayBuffer] is 16 bytes transaction key
//   - plainText[ArrayBuffer] is transaction request data, must not be empty.
//   - action is "request" or "response", anything else is an error.
//     Stricter than upstream (moov-io).
//
// Return Params:
//   - result[ArrayBuffer] - 8 bytes mac, use the first 4 bytes as the ANSI X9.24-1 mac
//   - err
func (m *module) GenerateMac(currentKeyArg, plainTextArg sobek.Value, action string) (*sobek.ArrayBuffer, error) {
	currentKey, err := jsbytes.Required("currentKey", currentKeyArg)
	if err != nil {
		return nil, err
	}
	rawPlainText, err := jsbytes.Required("plainText", plainTextArg)
	if err != nil {
		return nil, err
	}

	if err := validateAction(action); err != nil {
		return nil, err
	}
	if err := validateKey("currentKey", currentKey); err != nil {
		return nil, err
	}

	if err := validatePlainText(rawPlainText); err != nil {
		return nil, err
	}

	mac, err := des.GenerateMac(currentKey, string(rawPlainText), action)
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

func validateAction(a string) error {
	if a != pkg.ActionRequest && a != pkg.ActionResponse {
		return fmt.Errorf(`action must be "request" or "response", got %q`, a)
	}
	return nil
}

func requireLen(name string, b []byte, n int) error {
	if len(b) != n {
		return fmt.Errorf("%s must be %d bytes, got %d", name, n, len(b))
	}
	return nil
}

// validateKey checks a 16 byte DUKPT key, moov-io pads or truncates other lengths.
func validateKey(name string, key []byte) error {
	return requireLen(name, key, desKeyLen)
}

// validateKsn checks a 10 byte key serial number, moov-io does not.
func validateKsn(ksn []byte) error {
	return requireLen("ksn", ksn, desKsnLen)
}

// maxTcOneBits is the most one-bits an ANSI X9.24-1 originator sets in the
// transaction counter.
const maxTcOneBits = 10

// validateTransactionCounter rejects a counter no originator produces. moov-io
// derives from it anyway, and for counter 0 hands back the initial key itself.
func validateTransactionCounter(ksn []byte) error {
	tc := pkg.GetDesTcFromKsn(ksn)
	if tc == 0 {
		return errors.New("ksn transaction counter must not be 0")
	}
	if n := bits.OnesCount32(tc); n > maxTcOneBits {
		return fmt.Errorf("ksn transaction counter must have at most %d one-bits, got %d", maxTcOneBits, n)
	}
	return nil
}

// validateCiphertext requires whole, non-empty CBC blocks. moov-io zero pads a
// ragged ciphertext before decrypting, so a truncated one yields garbage
// instead of an error.
func validateCiphertext(ct []byte) error {
	if len(ct) == 0 || len(ct)%desBlockLen != 0 {
		return fmt.Errorf("ciphertext length must be a non-zero multiple of %d bytes, got %d", desBlockLen, len(ct))
	}
	return nil
}

// validatePinBlock checks an 8 byte encrypted PIN block. moov-io rejects other
// lengths too, but without naming the argument or the actual length.
func validatePinBlock(ct []byte) error {
	return requireLen("ciphertextPin", ct, desBlockLen)
}

// validatePlainText rejects empty data. moov-io encrypts it into an empty
// ciphertext that decryptData refuses, and MACs it as the zero block.
func validatePlainText(pt []byte) error {
	if len(pt) == 0 {
		return errors.New("plainText must not be empty")
	}
	return nil
}
