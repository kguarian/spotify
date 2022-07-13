// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "github.com/jkravitz/mytrace"

func installTestHooks() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
}

func uninstallTestHooks() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
}

// forceCloseSockets must be called only from TestMain.
func forceCloseSockets() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
}

func enableSocketConnect() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
}

func disableSocketConnect(network string) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
}
