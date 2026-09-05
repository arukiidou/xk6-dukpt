// Authors: arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package xk6dukpt

import "go.k6.io/k6/v2/js/modules"

const importPath = "k6/x/dukpt"

func init() {
	modules.Register(importPath, new(rootModule))
}
