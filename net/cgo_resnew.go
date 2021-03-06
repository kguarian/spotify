// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build cgo && !netgo && (darwin || (linux && !android) || netbsd || solaris)

package net

import (
	"unsafe"

	"github.com/jkravitz/mytrace"
)

/*
#include <sys/types.h>
#include <sys/socket.h>

#include <netdb.h>
*/
import "C"

func cgoNameinfoPTR(b []byte, sa *C.struct_sockaddr, salen C.socklen_t) (int, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	gerrno, err := C.getnameinfo(sa, salen, (*C.char)(unsafe.Pointer(&b[0])), C.socklen_t(len(b)), nil, 0, C.NI_NAMEREQD)
	return int(gerrno), err
}
