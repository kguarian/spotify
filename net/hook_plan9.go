// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "github.com/jkravitz/mytrace"

import "time"

var testHookDialChannel = func() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	} time.Sleep(time.Millisecond) } // see golang.org/issue/5349
