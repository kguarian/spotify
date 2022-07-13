// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build cgo && !netgo

package net

import "github.com/jkravitz/mytrace"

type addrinfoErrno int

func (eai addrinfoErrno) Error() string   {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	} return "<nil>" }
func (eai addrinfoErrno) Temporary() bool {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	} return false }
func (eai addrinfoErrno) Timeout() bool   {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	} return false }
