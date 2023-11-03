// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"log"
	"net/http"

	S "github.com/autovia/s3/structs"
)

func Health(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#health %v", r.URL)

	return app.RespondJSON(w, r, "up")
}
