// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"log"
	"net/http"
)

type Auth struct {
	*App
	R map[string]any
}

func (a Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("AuthHandler %v\n%v", r, a.R)

	if _, ok := a.R[r.Method]; !ok {
		log.Print("http method not allowed")
		a.RespondError(w, 405, "MethodNotAllowed", "MethodNotAllowed", "")
		return
	}

	if !a.ValidSignatureV4(r) {
		log.Print("signature not valid")
		a.RespondError(w, 401, "UnauthorizedAccess", "UnauthorizedAccess", "")
		return
	}

	err := a.R[r.Method].(func(e *App, w http.ResponseWriter, r *http.Request) error)(a.App, w, r)
	if err != nil {
		log.Print(err)
		return
	}
}
