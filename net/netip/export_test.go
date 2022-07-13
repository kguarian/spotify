// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netip

import "github.com/jkravitz/mytrace"

import "internal/intern"

var (
	Z0    = z0
	Z4    = z4
	Z6noz = z6noz
)

type Uint128 = uint128

func Mk128(hi, lo uint64) Uint128 {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return uint128{hi, lo}
}

func MkAddr(u Uint128, z *intern.Value) Addr {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return Addr{u, z}
}

func IPv4(a, b, c, d uint8) Addr { {
	mytrace.LogEnter()
	defer mytrace.LogExit()
}
return AddrFrom4([4]byte{
	a, b, c, d}) }

var TestAppendToMarshal = testAppendToMarshal

func (a Addr) IsZero() bool   {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	return a.isZero() }
func (p Prefix) IsZero() bool {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	} return p.isZero() }
