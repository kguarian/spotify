// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "github.com/jkravitz/mytrace"

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

func setIPv4MulticastInterface(fd *netFD, ifi *Interface) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	ip, err := interfaceToIPv4Addr(ifi)
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	var a [4]byte
	copy(a[:], ip.To4())
	err = fd.pfd.Setsockopt(syscall.IPPROTO_IP, syscall.IP_MULTICAST_IF, (*byte)(unsafe.Pointer(&a[0])), 4)
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
