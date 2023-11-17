// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (app *App) ValidSignatureV4(r *http.Request) (bool, *http.Request) {
	//log.Printf("---%s---", r.Header.Get("Authorization"))

	headers := authorizationHeader(r.Header, r.Host, r.Header.Get("Authorization"))
	if headers == nil {
		return false, nil
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	query := r.URL.Query()
	canonicalRequest := canonicalRequest(r.Method, r.URL, query.Encode(), bodyBytes, headers)
	if canonicalRequest == "" {
		return false, nil
	}

	stringToSign := stringToSign(canonicalRequest, *app.AccessKey, headers)
	if stringToSign == "" {
		return false, nil
	}

	signature := signingKeySignature(*app.SecretKey, stringToSign, headers)
	if signature == "" {
		return false, nil
	}

	//log.Print(signature)
	//log.Print(headers["Signature"])

	return signature == headers["Signature"], r
}

func authorizationHeader(header http.Header, host string, req string) map[string]string {
	authHeader, found := strings.CutPrefix(req, "AWS4-HMAC-SHA256")
	if !found {
		return nil
	}

	authHeaderSplit := strings.Split(strings.TrimSpace(authHeader), ",")
	if len(authHeaderSplit) != 3 {
		return nil
	}

	headers := make(map[string]string)
	for _, h := range authHeaderSplit {
		tuple := strings.Split(strings.TrimSpace(h), "=")
		if len(tuple) != 2 {
			return nil
		}
		headers[tuple[0]] = tuple[1]
	}

	headers["host"] = host
	for k, v := range header {
		headers[strings.ToLower(k)] = strings.Join(v, ",")
	}

	return headers
}

func canonicalRequest(method string, requestURI *url.URL, rawQuery string, body []byte, headers map[string]string) string {
	signedHeaders := strings.Split(headers["SignedHeaders"], ";")

	canonicalRequest := method + "\n"                                    // <HTTPMethod>
	canonicalRequest += requestURI.EscapedPath() + "\n"                  // <CanonicalURI>
	canonicalRequest += strings.Replace(rawQuery, "+", "%20", -1) + "\n" // <CanonicalQueryString>

	for _, v := range signedHeaders {
		canonicalRequest += v + ":" + strings.Join(strings.Fields(headers[v]), " ") + "\n" // <CanonicalHeaders>
	}
	canonicalRequest += "\n"
	canonicalRequest += headers["SignedHeaders"] + "\n" // <SignedHeaders>
	canonicalRequest += HexSHA256Hash(body)

	//log.Print(canonicalRequest)

	return canonicalRequest
}

func stringToSign(canonicalRequest string, accessKey string, headers map[string]string) string {
	scope, ok := strings.CutPrefix(headers["Credential"], accessKey+"/")
	if !ok {
		return ""
	}

	stringToSign := "AWS4-HMAC-SHA256\n"
	stringToSign += headers["x-amz-date"] + "\n"
	stringToSign += scope + "\n"
	stringToSign += HexSHA256Hash([]byte(canonicalRequest))

	//log.Print(stringToSign)

	return stringToSign
}

func signingKeySignature(secret string, stringToSign string, headers map[string]string) string {
	signPayload := strings.Split(headers["Credential"], "/")
	if len(signPayload) != 5 {
		return ""
	}

	dateKey := HmacSHA256([]byte("AWS4"+secret), []byte(signPayload[1]))
	dateRegionKey := HmacSHA256(dateKey, []byte(signPayload[2]))
	dateRegionServiceKey := HmacSHA256(dateRegionKey, []byte(signPayload[3]))
	signingKey := HmacSHA256(dateRegionServiceKey, []byte(signPayload[4]))

	return hex.EncodeToString(HmacSHA256(signingKey, []byte(stringToSign)))
}

func HexSHA256Hash(b []byte) string {
	bs := sha256.Sum256([]byte(b))
	return hex.EncodeToString(bs[:])
}

func HmacSHA256(key []byte, data []byte) []byte {
	hmac := hmac.New(sha256.New, key)
	hmac.Write(data)
	return hmac.Sum(nil)
}
