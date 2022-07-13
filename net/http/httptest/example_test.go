// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptest_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

func ExampleResponseRecorder() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		{
			mytrace.LogEnter()
			defer mytrace.LogExit()
		}
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Output:
	// 200
	// text/html; charset=utf-8
	// <html><body>Hello World!</body></html>
}

func ExampleServer() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	ts := httptest.NewServer(http.Handlerfunc(func(w http.ResponseWriter, r *http.Request) {
		{
			mytrace.LogEnter()
			defer mytrace.LogExit()
		}
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)
	// Output: Hello, client
}

func ExampleServer_hTTP2() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	ts := httptest.NewUnstartedServer(http.Handlerfunc(func(w http.ResponseWriter, r *http.Request) {
		{
			mytrace.LogEnter()
			defer mytrace.LogExit()
		}
		fmt.Fprintf(w, "Hello, %s", r.Proto)
	}))
	ts.EnableHTTP2 = true
	ts.StartTLS()
	defer ts.Close()

	res, err := ts.Client().Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", greeting)

	// Output: Hello, HTTP/2.0
}

func ExampleNewTLSServer() {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	ts := httptest.NewTLSServer(http.Handlerfunc(func(w http.ResponseWriter, r *http.Request) {
		{
			mytrace.LogEnter()
			defer mytrace.LogExit()
		}
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	client := ts.Client()
	res, err := client.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)
	// Output: Hello, client
}