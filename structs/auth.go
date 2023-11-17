// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"errors"
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
		a.RespondError(w, 405, "MethodNotAllowed", errors.New("MethodNotAllowed"), "")
		return
	}

	valid, req := a.ValidSignatureV4(r)
	if !valid {
		log.Print("signature not valid")
		a.RespondError(w, 401, "UnauthorizedAccess", errors.New("UnauthorizedAccess"), "")
		return
	}

	err := a.R[req.Method].(func(e *App, w http.ResponseWriter, r *http.Request) error)(a.App, w, req)
	if err != nil {
		log.Print(err)
		return
	}
}
