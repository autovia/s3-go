// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
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

func (app *App) RespondXML(w http.ResponseWriter, code int, payload any) error {

	out, _ := xml.MarshalIndent(payload, " ", "  ")
	//fmt.Println(string(out))

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	w.Write([]byte(out))

	return nil
}

func (app *App) Respond(w http.ResponseWriter, code int, headers map[string]string, body []byte) error {
	if len(headers) > 0 {
		for k, v := range headers {
			w.Header().Set(k, v)
			log.Print(k, v)
		}
	}

	w.WriteHeader(code)

	// if len(body) > 0 {
	// 	w.Write(body)
	// }

	return nil
}

func (app *App) RespondFile(w http.ResponseWriter, code int, headers map[string]string, file *os.File) error {
	if len(headers) > 0 {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}

	w.WriteHeader(code)
	defer file.Close()
	io.Copy(w, file)

	return nil
}
