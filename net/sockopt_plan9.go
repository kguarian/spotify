// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "github.com/jkravitz/mytrace"

import "syscall"

func setKeepAlive(fd *netFD, keepalive bool) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	if keepalive {
		_, e := fd.ctl.WriteAt([]byte("keepalive"), 0)
		return e
	}
	return nil
}

func setLinger(fd *netFD, sec int) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.EPLAN9
}
