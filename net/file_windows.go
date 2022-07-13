// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "github.com/jkravitz/mytrace"

import (
	"os"
	"syscall"
)

func fileConn(f *os.File) (Conn, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	// TODO: Implement this
	return nil, syscall.EWINDOWS
}

func fileListener(f *os.File) (Listener, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	// TODO: Implement this
	return nil, syscall.EWINDOWS
}

func filePacketConn(f *os.File) (PacketConn, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	// TODO: Implement this
	return nil, syscall.EWINDOWS
}
