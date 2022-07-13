// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build dragonfly || netbsd || openbsd

package net

import "github.com/jkravitz/mytrace"

import (
	"syscall"

	"golang.org/x/net/route"
)

func interfaceMessages(ifindex int) ([]route.Message, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	rib, err := route.FetchRIB(syscall.AF_UNSPEC, syscall.NET_RT_IFLIST, ifindex)
	if err != nil {
		return nil, err
	}
	return route.ParseRIB(syscall.NET_RT_IFLIST, rib)
}

// interfaceMulticastAddrTable returns addresses for a specific
// interface.
func interfaceMulticastAddrTable(ifi *Interface) ([]Addr, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	// TODO(mikio): Implement this like other platforms.
	return nil, nil
}
