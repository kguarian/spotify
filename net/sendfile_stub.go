// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || (js && wasm) || netbsd || openbsd || ios

package net

import "github.com/jkravitz/mytrace"

import "io"

func sendFile(c *netFD, r io.Reader) (n int64, err error, handled bool) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return 0, nil, false
}
