// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"context"
	"log"
	"net/http"
)

type Public struct {
	*App
	R map[string]any
}

func (p Public) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("PublicHandler %v", r.URL)

	if _, ok := p.R[r.Method]; !ok {
		log.Print("http method not allowed")
		//HandleError(w, fmt.Errorf("http method not allowed"))
		return
	}

	err := p.R[r.Method].(func(e *App, w http.ResponseWriter, r *http.Request) error)(p.App, w, r)
	if err != nil {
		log.Print(err)
		//HandleError(w, err)
		return
	}
}

type Auth struct {
	*App
	R map[string]any
}

func (a Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("AuthHandler %v\n%v", r, a.R)

	if _, ok := a.R[r.Method]; !ok {
		log.Print("http method not allowed")
		//HandleError(w, fmt.Errorf("http method not allowed"))
		return
	}

	if !a.ValidSignatureV4(r) {
		log.Print("signature not valid")
		//HandleError(w, fmt.Errorf("http method not allowed"))
		return
	}

	token := "test"
	ctx := context.WithValue(r.Context(), "csrf", token)
	newReq := r.WithContext(ctx)
	log.Printf("token %s", token)

	err := a.R[r.Method].(func(e *App, w http.ResponseWriter, r *http.Request) error)(a.App, w, newReq)
	if err != nil {
		log.Print(err)
		//HandleError(w, err)
		return
	}
}
