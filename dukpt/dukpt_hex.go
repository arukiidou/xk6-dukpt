// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"encoding/hex"

	"github.com/moov-io/dukpt/pkg"
)

// [pkg.GetDesTcFromKsn] port from moov-io, as a Hex variant
//
// ksn hex string - 10 bytes key serial number.
//
// Return Params:
//   - result is 21 bits transaction counter
//   - err
//
// KSN length is not validated.
func GetDesTcFromKsnAsHex(ksn string) (uint32, error) {
	rawKsn, err := hex.DecodeString(ksn)
	if err != nil {
		return 0, err
	}
	return pkg.GetDesTcFromKsn(rawKsn), nil
}
