// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import "github.com/moov-io/dukpt/pkg"

// [pkg.GetDesTcFromKsn] port from moov-io
//
// Params:
//   - ksn[ArrayBuffer] is 10 bytes key serial number
//
// Return Params:
//   - result - 21 bits transaction counter
//
// KSN length is not validated.
func GetDesTcFromKsn(ksn []byte) uint32 {
	return pkg.GetDesTcFromKsn(ksn)
}
