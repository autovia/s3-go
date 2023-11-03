// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/autovia/s3/handlers"
	S "github.com/autovia/s3/structs"
)

func InitRoutes(app *S.App) {
	// Public handlers
	app.Router.Handle("/health", S.Public{App: app, R: map[string]any{
		"GET": handlers.Health,
	}})

	app.Router.Handle("/", S.Auth{App: app, R: map[string]any{
		"GET": handlers.ListBuckets,
		"PUT": handlers.CreateBucket,
	}})
}
