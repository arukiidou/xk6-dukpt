// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

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
