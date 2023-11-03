// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
)

func (app *App) RespondJSON(w http.ResponseWriter, r *http.Request, payload any) error {
	out, _ := json.Marshal(payload)
	//fmt.Println(string(out))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(out))

	log.Printf("%v", r)

	return nil
}

func (app *App) RespondXML(w http.ResponseWriter, r *http.Request, payload any) error {

	out, _ := xml.MarshalIndent(payload, " ", "  ")
	//fmt.Println(string(out))

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(out))

	return nil
}
