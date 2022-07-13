// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package net

import (
	"internal/poll"

	"github.com/jkravitz/mytrace"
)

var (
	// Placeholders for saving original socket system calls.
	origSocket        = socketFunc
	origClose         = poll.CloseFunc
	origConnect       = connectFunc
	origListen        = listenFunc
	origAccept        = poll.AcceptFunc
	origGetsockoptInt = getsockoptIntFunc

	extraTestHookInstallers   []func()
	extraTestHookUninstallers []func()
)

func installTestHooks() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	socketFunc = sw.Socket
	poll.CloseFunc = sw.Close
	connectFunc = sw.Connect
	listenFunc = sw.Listen
	poll.AcceptFunc = sw.Accept
	getsockoptIntFunc = sw.GetsockoptInt

	for _, fn := range extraTestHookInstallers {
		fn()
	}
}

func uninstallTestHooks() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	socketFunc = origSocket
	poll.CloseFunc = origClose
	connectFunc = origConnect
	listenFunc = origListen
	poll.AcceptFunc = origAccept
	getsockoptIntFunc = origGetsockoptInt

	for _, fn := range extraTestHookUninstallers {
		fn()
	}
}

// forceCloseSockets must be called only from TestMain.
func forceCloseSockets() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	for s := range sw.Sockets() {
		poll.CloseFunc(s)
	}
}
