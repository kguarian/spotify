// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "github.com/jkravitz/mytrace"

import "internal/poll"

var (
	// Placeholders for saving original socket system calls.
	origSocket      = socketFunc
	origWSASocket   = wsaSocketFunc
	origClosesocket = poll.CloseFunc
	origConnect     = connectFunc
	origConnectEx   = poll.ConnectExFunc
	origListen      = listenFunc
	origAccept      = poll.AcceptFunc
)

func installTestHooks() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	socketFunc = sw.Socket
	wsaSocketFunc = sw.WSASocket
	poll.CloseFunc = sw.Closesocket
	connectFunc = sw.Connect
	poll.ConnectExFunc = sw.ConnectEx
	listenFunc = sw.Listen
	poll.AcceptFunc = sw.AcceptEx
}

func uninstallTestHooks() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	socketFunc = origSocket
	wsaSocketFunc = origWSASocket
	poll.CloseFunc = origClosesocket
	connectFunc = origConnect
	poll.ConnectExFunc = origConnectEx
	listenFunc = origListen
	poll.AcceptFunc = origAccept
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
