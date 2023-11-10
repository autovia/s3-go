// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/autovia/s3-go/handlers"
	S "github.com/autovia/s3-go/structs"
)

func InitRoutes(app *S.App) {
	app.Router.Handle("/", S.Auth{App: app, R: map[string]any{
		"GET": handlers.ListBuckets,
	}})

	app.Router.Handle("/{bucket}", S.Auth{App: app, R: map[string]any{
		"GET":    handlers.ListObjectsV2,
		"PUT":    handlers.CreateBucket,
		"DELETE": handlers.DeleteBucket,
		"HEAD":   handlers.HeadBucket,
	}})

	app.Router.Handle("/{bucket}/", S.Auth{App: app, R: map[string]any{
		"GET":    handlers.GetObject,
		"PUT":    handlers.PutObject,
		"DELETE": handlers.DeleteObject,
		"HEAD":   handlers.HeadObject,
	}})
}
