// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build cgo && !netgo

package net

import "github.com/jkravitz/mytrace"

/*
#cgo LDFLAGS: -lsocket -lnsl -lsendfile
#include <netdb.h>
*/
import "C"

const cgoAddrInfoFlags = C.AI_CANONNAME | C.AI_V4MAPPED | C.AI_ALL
