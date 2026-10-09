// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"encoding/base64"
)

// [pkg.GetDesTcFromKsn] port from moov-io
//
// ksn base64 string - 10 bytes key serial number.
//
// Return Params:
//   - result is 21 bits transaction counter
//   - err
//
// The KSN must be at least 4 bytes, moov-io returns 0 for shorter ones.
func GetDesTcFromKsnAsBase64(ksn string) (uint32, error) {
	rawKsn, err := base64.StdEncoding.DecodeString(ksn)
	if err != nil {
		return 0, err
	}
	return getDesTcFromKsn(rawKsn)
}

// [pkg.GenerateNextDesKsn] port from moov-io
//
// ksn base64 string - 10 bytes key serial number.
//
// Return Params:
//   - result is base64 string - 10 bytes next key serial number
//   - err
func GenerateNextDesKsnAsBase64(ksn string) (string, error) {
	rawKsn, err := base64.StdEncoding.DecodeString(ksn)
	if err != nil {
		return "", err
	}

	next, err := generateNextDesKsn(rawKsn)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(next), nil
}
