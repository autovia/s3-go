// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestValidSignatureV4(t *testing.T) {
	headerStr := "AWS4-HMAC-SHA256 Credential=user/20130524/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=784a174cd7fea56e95c9b5ca2eb6aa5e78532e9b58f58780aa65309526083d0a"

	expectedAuthorizationHeader := map[string]string{
		"Credential":           "user/20130524/us-east-1/s3/aws4_request",
		"SignedHeaders":        "host;x-amz-content-sha256;x-amz-date",
		"Signature":            "784a174cd7fea56e95c9b5ca2eb6aa5e78532e9b58f58780aa65309526083d0a",
		"host":                 "examplebucket.s3.amazonaws.com",
		"x-amz-content-sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"x-amz-date":           "20130524T000000Z",
	}

	header := map[string][]string{
		"X-Amz-Content-Sha256": {"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		"X-Amz-Date":           {"20130524T000000Z"},
	}

	t.Run("authorizationHeader", func(t *testing.T) {
		result := authorizationHeader(header, "examplebucket.s3.amazonaws.com", headerStr)
		if !reflect.DeepEqual(expectedAuthorizationHeader, result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedAuthorizationHeader)
		}
	})

	expectedCanonicalRequest := `GET
/
lifecycle=
host:examplebucket.s3.amazonaws.com
x-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
x-amz-date:20130524T000000Z

host;x-amz-content-sha256;x-amz-date
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`

	t.Run("canonicalRequest", func(t *testing.T) {
		headers := authorizationHeader(header, "examplebucket.s3.amazonaws.com", headerStr)

		result := canonicalRequest("GET", "/", "lifecycle=", io.NopCloser(bytes.NewReader([]byte(""))), headers)
		if expectedCanonicalRequest != result {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedCanonicalRequest)
		}
	})

	expectedStringToSign := `AWS4-HMAC-SHA256
20130524T000000Z
20130524/us-east-1/s3/aws4_request
9766c798316ff2757b517bc739a67f6213b4ab36dd5da2f94eaebf79c77395ca`

	t.Run("stringToSign", func(t *testing.T) {
		result := stringToSign(expectedCanonicalRequest, "user", expectedAuthorizationHeader)
		if expectedStringToSign != result {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedStringToSign)
		}
	})

	expectedSigningKey := "fea454ca298b7da1c68078a5d1bdbfbbe0d65c699e0f91ac7a200a0136783543"

	t.Run("signingKeySignature", func(t *testing.T) {
		result := signingKeySignature("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", expectedStringToSign, expectedAuthorizationHeader)
		if expectedSigningKey != string(result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedSigningKey)
		}
	})

}
