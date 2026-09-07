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
