// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "github.com/jkravitz/mytrace"

import (
	"syscall"
	"time"
)

func setKeepAlivePeriod(fd *netFD, d time.Duration) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	// OpenBSD has no user-settable per-socket TCP keepalive
	// options.
	return syscall.ENOPROTOOPT
}
