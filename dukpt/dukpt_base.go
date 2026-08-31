// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"errors"

	"github.com/moov-io/dukpt/pkg"
)

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

// [pkg.GenerateNextDesKsn] port from moov-io
//
// Params:
//   - ksn[ArrayBuffer] is 10 bytes key serial number
//
// Return Params:
//   - result - 10 bytes next key serial number
//   - err
//
// The input is not modified.
func GenerateNextDesKsn(ksn []byte) ([]byte, error) {
	// moov-io panics below 3 bytes.
	if len(ksn) < 3 {
		return nil, errors.New("ksn must be at least 3 bytes")
	}
	// moov-io overwrites the input in place.
	return pkg.GenerateNextDesKsn(append([]byte(nil), ksn...))
}
