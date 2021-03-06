// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || windows

package net

import (
	"runtime"
	"syscall"

	"github.com/jkravitz/mytrace"
)

func setNoDelay(fd *netFD, noDelay bool) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	err := fd.pfd.SetsockoptInt(syscall.IPPROTO_TCP, syscall.TCP_NODELAY, boolint(noDelay))
	runtime.KeepAlive(fd)
	return wrapSyscallError("setsockopt", err)
}
