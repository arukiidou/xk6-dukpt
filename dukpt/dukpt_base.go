// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package dukpt

import (
	"fmt"

	"github.com/grafana/sobek"
	"github.com/moov-io/dukpt/pkg"
)

// newArrayBuffer hands raw bytes back as JS ArrayBuffer.
func (m *module) newArrayBuffer(b []byte) *sobek.ArrayBuffer {
	ab := m.vu.Runtime().NewArrayBuffer(b)
	return &ab
}

// [pkg.GetDesTcFromKsn] port from moov-io
//
// Params:
//   - ksn[ArrayBuffer] is 10 bytes key serial number
//
// Return Params:
//   - result - 21 bits transaction counter
//   - err
//
// The KSN must be at least 4 bytes, moov-io returns 0 for shorter ones.
func GetDesTcFromKsn(ksn []byte) (uint32, error) {
	return getDesTcFromKsn(ksn)
}

// minTcKsnLen is the shortest KSN whose transaction counter moov-io reads.
const minTcKsnLen = 4

// getDesTcFromKsn is the raw form the hex and base64 variants share.
func getDesTcFromKsn(ksn []byte) (uint32, error) {
	// moov-io silently reports counter 0 below 4 bytes.
	if len(ksn) < minTcKsnLen {
		return 0, fmt.Errorf("ksn must be at least %d bytes, got %d", minTcKsnLen, len(ksn))
	}
	return pkg.GetDesTcFromKsn(ksn), nil
}

// [pkg.GenerateNextDesKsn] port from moov-io
//
// Params:
//   - ksn[ArrayBuffer] is 10 bytes key serial number
//
// Return Params:
//   - result[ArrayBuffer] - 10 bytes next key serial number
//   - err
//
// The input is not modified.
func (m *module) GenerateNextDesKsn(ksn []byte) (*sobek.ArrayBuffer, error) {
	next, err := generateNextDesKsn(ksn)
	if err != nil {
		return nil, err
	}
	return m.newArrayBuffer(next), nil
}

// generateNextDesKsn is the raw form the hex and base64 variants encode.
func generateNextDesKsn(ksn []byte) ([]byte, error) {
	// moov-io panics below 3 bytes, and restarts the counter at 1 below 4.
	if len(ksn) < minTcKsnLen {
		return nil, fmt.Errorf("ksn must be at least %d bytes, got %d", minTcKsnLen, len(ksn))
	}
	// moov-io overwrites the input in place.
	return pkg.GenerateNextDesKsn(append([]byte(nil), ksn...))
}
