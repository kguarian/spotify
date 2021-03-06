// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !js && !plan9

package socktest_test

import (
	"net/internal/socktest"
	"os"
	"sync"
	"syscall"
	"testing"
)

var sw socktest.Switch

func TestMain(m *testing.M) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	installTestHooks()

	st := m.Run()

	for s := range sw.Sockets() {
		closeFunc(s)
	}
	uninstallTestHooks()
	os.Exit(st)
}

func TestSwitch(t *testing.T) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	const N = 10
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
			defer wg.Done()
			for _, family := range []int{syscall.AF_INET, syscall.AF_INET6} {
				socketFunc(family, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
			}
		}()
	}
	wg.Wait()
}

func TestSocket(t *testing.T) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	for _, f := range []socktest.Filter{
		func(st *socktest.Status) (socktest.AfterFilter, error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	} return nil, nil },
		nil,
	} {
		sw.Set(socktest.FilterSocket, f)
		for _, family := range []int{syscall.AF_INET, syscall.AF_INET6} {
			socketFunc(family, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
		}
	}
}
