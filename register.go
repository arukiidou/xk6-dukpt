// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

// Package xk6dukpt is a k6 extension for DUKPT (Derived Unique Key Per Transaction) cryptography operations.
package xk6dukpt

import "go.k6.io/k6/v2/js/modules"

const importPath = "k6/x/dukpt"

func init() {
	modules.Register(importPath, new(rootModule))
}
