// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func (app *App) RespondXML(w http.ResponseWriter, code int, payload any) error {
	out, _ := xml.MarshalIndent(payload, " ", "  ")
	fmt.Println(string(out))

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	w.Write([]byte(out))

	return nil
}

func (app *App) Respond(w http.ResponseWriter, code int, headers map[string]string, body []byte) error {
	if len(headers) > 0 {
		for k, v := range headers {
			w.Header().Set(k, v)
			log.Println("metadata >>", k, v)
		}
	}

	w.WriteHeader(code)
	if len(body) > 0 {
		w.Write(body)
	}

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

type Error struct {
	XMLName   xml.Name `xml:"Error"`
	Code      string   `xml:"Code"`
	Message   string   `xml:"Message"`
	Resource  string   `xml:"Resource"`
	RequestId string   `xml:"RequestId"`
}

func (app *App) RespondError(w http.ResponseWriter, httpcode int, awscode string, err error, resource string) error {
	e := Error{
		Code:     awscode,
		Message:  awscode,
		Resource: resource,
	}

	log.Print(">>>", err)

	out, _ := xml.MarshalIndent(e, " ", "  ")
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(httpcode)
	w.Write([]byte(out))

	return nil
}
