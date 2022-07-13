// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socktest_test

import "syscall"

var (
	socketFunc func(int, int, int) (syscall.Handle, error)
	closeFunc  func(syscall.Handle) error
)

func installTestHooks() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	socketFunc = sw.Socket
	closeFunc = sw.Closesocket
}

func uninstallTestHooks() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	socketFunc = syscall.Socket
	closeFunc = syscall.Closesocket
}
