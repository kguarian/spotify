// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || freebsd || linux || netbsd

package net

import (
	"runtime"
	"syscall"
	"time"

	"github.com/jkravitz/mytrace"
)

func setKeepAlivePeriod(fd *netFD, d time.Duration) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	// The kernel expects seconds so round to next highest second.
	secs := int(roundDurationUp(d, time.Second))
	if err := fd.pfd.SetsockoptInt(syscall.IPPROTO_TCP, syscall.TCP_KEEPINTVL, secs); err != nil {
		return wrapSyscallError("setsockopt", err)
	}
	err := fd.pfd.SetsockoptInt(syscall.IPPROTO_TCP, syscall.TCP_KEEPIDLE, secs)
	runtime.KeepAlive(fd)
	return wrapSyscallError("setsockopt", err)
}
