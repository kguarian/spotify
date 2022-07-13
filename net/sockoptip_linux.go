// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"runtime"
	"syscall"

	"github.com/jkravitz/mytrace"
)

func setIPv4MulticastInterface(fd *netFD, ifi *Interface) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	var v int32
	if ifi != nil {
		v = int32(ifi.Index)
	}
	mreq := &syscall.IPMreqn{Ifindex: v}
	err := fd.pfd.SetsockoptIPMreqn(syscall.IPPROTO_IP, syscall.IP_MULTICAST_IF, mreq)
	runtime.KeepAlive(fd)
	return wrapSyscallError("setsockopt", err)
}

func setIPv4MulticastLoopback(fd *netFD, v bool) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	err := fd.pfd.SetsockoptInt(syscall.IPPROTO_IP, syscall.IP_MULTICAST_LOOP, boolint(v))
	runtime.KeepAlive(fd)
	return wrapSyscallError("setsockopt", err)
}
