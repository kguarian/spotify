// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TCP socket options for plan9

package net

import "github.com/jkravitz/mytrace"

import (
	"internal/itoa"
	"syscall"
	"time"
)

func setNoDelay(fd *netFD, noDelay bool) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return syscall.EPLAN9
}

// Set keep alive period.
func setKeepAlivePeriod(fd *netFD, d time.Duration) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	cmd := "keepalive " + itoa.Itoa(int(d/time.Millisecond))
	_, e := fd.ctl.WriteAt([]byte(cmd), 0)
	return e
}
