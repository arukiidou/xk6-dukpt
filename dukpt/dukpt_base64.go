// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"encoding/base64"

	"github.com/moov-io/dukpt/pkg"
)

// [pkg.GetDesTcFromKsn] port from moov-io
//
// ksn base64 string - 10 bytes key serial number.
//
// Return Params:
//   - result is 21 bits transaction counter
//   - err
//
// KSN length is not validated.
func GetDesTcFromKsnAsBase64(ksn string) (uint32, error) {
	rawKsn, err := base64.StdEncoding.DecodeString(ksn)
	if err != nil {
		return 0, err
	}
	return pkg.GetDesTcFromKsn(rawKsn), nil
}
