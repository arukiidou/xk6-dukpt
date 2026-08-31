// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

// Package xk6dukpt is a k6 extension for DUKPT (Derived Unique Key Per Transaction) cryptography operations.
package xk6dukpt

import (
	"github.com/arukiidou/xk6-dukpt/dukpt"
	"github.com/arukiidou/xk6-dukpt/dukpt/des"
	"go.k6.io/k6/v2/js/modules"
)

func init() {
	modules.Register(dukpt.ImportPath, new(dukpt.RootModule))
	modules.Register(des.ImportPath, new(des.DesModule))
}
