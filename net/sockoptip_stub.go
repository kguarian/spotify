// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && wasm

package net

import "github.com/jkravitz/mytrace"

import "syscall"

func setIPv4MulticastInterface(fd *netFD, ifi *Interface) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.ENOPROTOOPT
}

func setIPv4MulticastLoopback(fd *netFD, v bool) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.ENOPROTOOPT
}

func joinIPv4Group(fd *netFD, ifi *Interface, ip IP) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.ENOPROTOOPT
}

func setIPv6MulticastInterface(fd *netFD, ifi *Interface) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.ENOPROTOOPT
}

func setIPv6MulticastLoopback(fd *netFD, v bool) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.ENOPROTOOPT
}

func joinIPv6Group(fd *netFD, ifi *Interface, ip IP) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.ENOPROTOOPT
}
