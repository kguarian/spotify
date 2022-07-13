// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !js && !plan9 && !windows

package socktest_test

import "syscall"

var (
	socketFunc func(int, int, int) (int, error)
	closeFunc  func(int) error
)

func installTestHooks() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	socketFunc = sw.Socket
	closeFunc = sw.Close
}

func uninstallTestHooks() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	socketFunc = syscall.Socket
	closeFunc = syscall.Close
}
