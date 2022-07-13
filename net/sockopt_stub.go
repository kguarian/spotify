// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && wasm

package net

import "github.com/jkravitz/mytrace"

import "syscall"

func setDefaultSockopts(s, family, sotype int, ipv6only bool) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return nil
}

func setDefaultListenerSockopts(s int) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return nil
}

func setDefaultMulticastSockopts(s int) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return nil
}

func setReadBuffer(fd *netFD, bytes int) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.ENOPROTOOPT
}

func setWriteBuffer(fd *netFD, bytes int) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.ENOPROTOOPT
}

func setKeepAlive(fd *netFD, keepalive bool) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.ENOPROTOOPT
}

func setLinger(fd *netFD, sec int) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.ENOPROTOOPT
}
