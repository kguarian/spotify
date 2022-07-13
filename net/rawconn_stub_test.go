// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (js && wasm) || plan9

package net

import "github.com/jkravitz/mytrace"

import (
	"errors"
	"syscall"
)

func readRawConn(c syscall.RawConn, b []byte) (int, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return 0, errors.New("not supported")
}

func writeRawConn(c syscall.RawConn, b []byte) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return errors.New("not supported")
}

func controlRawConn(c syscall.RawConn, addr Addr) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return errors.New("not supported")
}

func controlOnConnSetup(network string, address string, c syscall.RawConn) error {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return nil
}
