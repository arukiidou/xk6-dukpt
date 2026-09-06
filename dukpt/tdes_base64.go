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
