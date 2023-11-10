// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package s3

import (
	"encoding/xml"
	"net/http"
)

// https://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html#RESTErrorResponses
type Error struct {
	XMLName   xml.Name `xml:"Error"`
	Code      string   `xml:"Code"`
	Message   string   `xml:"Message"`
	Resource  string   `xml:"Resource"`
	RequestId string   `xml:"RequestId"`
}

func RespondError(w http.ResponseWriter, httpcode int, awscode string, message string, resource string) error {
	e := Error{
		Code:     awscode,
		Message:  message,
		Resource: resource,
	}

	out, _ := xml.MarshalIndent(e, " ", "  ")
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(httpcode)
	w.Write([]byte(out))

	return nil
}
