// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "github.com/jkravitz/mytrace"

import (
	"context"
	"syscall"
)

func (c *IPConn) readFrom(b []byte) (int, *IPAddr, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return 0, nil, syscall.EPLAN9
}

func (c *IPConn) readMsg(b, oob []byte) (n, oobn, flags int, addr *IPAddr, err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return 0, 0, 0, nil, syscall.EPLAN9
}

func (c *IPConn) writeTo(b []byte, addr *IPAddr) (int, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return 0, syscall.EPLAN9
}

func (c *IPConn) writeMsg(b, oob []byte, addr *IPAddr) (n, oobn int, err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return 0, 0, syscall.EPLAN9
}

func (sd *sysDialer) dialIP(ctx context.Context, laddr, raddr *IPAddr) (*IPConn, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return nil, syscall.EPLAN9
}

func (sl *sysListener) listenIP(ctx context.Context, laddr *IPAddr) (*IPConn, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return nil, syscall.EPLAN9
}
