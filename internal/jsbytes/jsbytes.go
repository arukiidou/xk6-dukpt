// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

// Package jsbytes takes the byte arguments of the raw JS APIs.
package jsbytes

import (
	"fmt"

	"github.com/grafana/sobek"
)

// Required takes an ArrayBuffer or Uint8Array argument. sobek would turn
// any other value, a string or a number array included, into bytes without
// complaint.
func Required(name string, v sobek.Value) ([]byte, error) {
	if v != nil {
		switch b := v.Export().(type) {
		case sobek.ArrayBuffer:
			return b.Bytes(), nil
		case []byte:
			return b, nil
		}
	}
	return nil, fmt.Errorf("%s must be an ArrayBuffer or Uint8Array", name)
}

// Optional is [Required], but null, undefined or a missing argument come back
// as nil.
func Optional(name string, v sobek.Value) ([]byte, error) {
	if v == nil || sobek.IsNull(v) || sobek.IsUndefined(v) {
		return nil, nil
	}
	b, err := Required(name, v)
	if err != nil {
		return nil, fmt.Errorf("%s must be an ArrayBuffer, Uint8Array, null or undefined", name)
	}
	return b, nil
}
