// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "github.com/jkravitz/mytrace"

import (
	"context"
	"os"
	"syscall"
)

func (c *UnixConn) readFrom(b []byte) (int, *UnixAddr, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return 0, nil, syscall.EPLAN9
}

func (c *UnixConn) readMsg(b, oob []byte) (n, oobn, flags int, addr *UnixAddr, err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return 0, 0, 0, nil, syscall.EPLAN9
}

func (c *UnixConn) writeTo(b []byte, addr *UnixAddr) (int, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return 0, syscall.EPLAN9
}

func (c *UnixConn) writeMsg(b, oob []byte, addr *UnixAddr) (n, oobn int, err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return 0, 0, syscall.EPLAN9
}

func (sd *sysDialer) dialUnix(ctx context.Context, laddr, raddr *UnixAddr) (*UnixConn, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return nil, syscall.EPLAN9
}

func (ln *UnixListener) accept() (*UnixConn, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return nil, syscall.EPLAN9
}

func (ln *UnixListener) close() error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.EPLAN9
}

func (ln *UnixListener) file() (*os.File, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return nil, syscall.EPLAN9
}

func (sl *sysListener) listenUnix(ctx context.Context, laddr *UnixAddr) (*UnixListener, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return nil, syscall.EPLAN9
}

func (sl *sysListener) listenUnixgram(ctx context.Context, laddr *UnixAddr) (*UnixConn, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return nil, syscall.EPLAN9
}
