// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9

package socktest

// Sockets maps a socket descriptor to the status of socket.
type Sockets map[int]Status

func familyString(family int) string {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	} return "<nil>" }

func typeString(sotype int) string {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	} return "<nil>" }

func protocolString(proto int) string {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	} return "<nil>" }
