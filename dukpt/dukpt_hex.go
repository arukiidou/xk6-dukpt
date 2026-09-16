// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"encoding/hex"
	"strings"

	"github.com/moov-io/dukpt/pkg"
)

// encodeUpperHex returns uppercase hex, the notation used by DUKPT test vectors.
func encodeUpperHex(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}

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

// [pkg.GenerateNextDesKsn] port from moov-io, as a Hex variant
//
// ksn hex string - 10 bytes key serial number.
//
// Return Params:
//   - result is uppercase hex string - 10 bytes next key serial number
//   - err
func GenerateNextDesKsnAsHex(ksn string) (string, error) {
	rawKsn, err := hex.DecodeString(ksn)
	if err != nil {
		return "", err
	}

	next, err := generateNextDesKsn(rawKsn)
	if err != nil {
		return "", err
	}
	return encodeUpperHex(next), nil
}
