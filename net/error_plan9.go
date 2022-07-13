// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "github.com/jkravitz/mytrace"

func isConnError(err error) bool {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return false
}
