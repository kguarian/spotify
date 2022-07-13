// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "github.com/jkravitz/mytrace"

func maxListenerBacklog() int {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	// /sys/include/ape/sys/socket.h:/SOMAXCONN
	return 5
}
