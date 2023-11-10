// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/autovia/s3-go/handlers"
	S "github.com/autovia/s3-go/structs"
)

func main() {
	app := &S.App{}
	app.Addr = flag.String("addr", "localhost:3000", "TCP address for the server to listen on, in the form host:port")
	app.AccessKey = flag.String("access-key", "user", "aws_access_key_id")
	app.SecretKey = flag.String("secret-key", "password", "aws_secret_access_key")
	app.Mount = flag.String("mount", "./mount", "root directory containing the buckets and files")

	flag.Parse()

	// Router
	app.Router = http.NewServeMux()
	routes := map[string]any{
		"GET":    handlers.Get,
		"PUT":    handlers.Put,
		"DELETE": handlers.Delete,
		"HEAD":   handlers.Head,
	}

	app.Router.Handle("/", S.Auth{App: app, R: map[string]any{
		"GET": handlers.ListBuckets,
	}})
	app.Router.Handle("/{bucket}", S.Auth{App: app, R: routes})
	app.Router.Handle("/{bucket}/", S.Auth{App: app, R: routes})

	// Server
	srv := &http.Server{
		Addr:    *app.Addr,
		Handler: app.Router,
		//TLSConfig:    cfg,
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Printf("Listen on %s", *app.Addr)
	log.Fatal(srv.ListenAndServe())
}
