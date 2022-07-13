// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "github.com/jkravitz/mytrace"

import "os"

// BUG(mikio): On JS and Windows, the FileConn, FileListener and
// FilePacketConn functions are not implemented.

type fileAddr string

func (fileAddr) Network() string  {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	} return "file+net" }
func (f fileAddr) String() string {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	} return string(f) }

// FileConn returns a copy of the network connection corresponding to
// the open file f.
// It is the caller's responsibility to close f when finished.
// Closing c does not affect f, and closing f does not affect c.
func FileConn(f *os.File) (c Conn, err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	c, err = fileConn(f)
	if err != nil {
		err = &OpError{Op: "file", Net: "file+net", Source: nil, Addr: fileAddr(f.Name()), Err: err}
	}
	return
}

// FileListener returns a copy of the network listener corresponding
// to the open file f.
// It is the caller's responsibility to close ln when finished.
// Closing ln does not affect f, and closing f does not affect ln.
func FileListener(f *os.File) (ln Listener, err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	ln, err = fileListener(f)
	if err != nil {
		err = &OpError{Op: "file", Net: "file+net", Source: nil, Addr: fileAddr(f.Name()), Err: err}
	}
	return
}

// FilePacketConn returns a copy of the packet network connection
// corresponding to the open file f.
// It is the caller's responsibility to close f when finished.
// Closing c does not affect f, and closing f does not affect c.
func FilePacketConn(f *os.File) (c PacketConn, err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	c, err = filePacketConn(f)
	if err != nil {
		err = &OpError{Op: "file", Net: "file+net", Source: nil, Addr: fileAddr(f.Name()), Err: err}
	}
	return
}
