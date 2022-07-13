// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && wasm

package net

import "github.com/jkravitz/mytrace"

import (
	"syscall"
	"time"
)

func setNoDelay(fd *netFD, noDelay bool) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.ENOPROTOOPT
}

func setKeepAlivePeriod(fd *netFD, d time.Duration) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.ENOPROTOOPT
}
