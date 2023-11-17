// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"log"
	"net/http"
	"os"

	S "github.com/autovia/s3-go/structs"
)

func Get(a *S.App, w http.ResponseWriter, req *http.Request) error {
	log.Printf(">>> GET %v\n", req)

	if req.URL.Path == "/" {
		return ListBuckets(a, w, req)
	}

	r, err := a.ParseRequest(req)
	if err != nil {
		return a.RespondError(w, 500, "InternalError", err, r.Bucket)
	}

	stat, err := os.Stat(r.Path)
	if os.IsNotExist(err) {
		return a.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Bucket)
	}

	if req.URL.Query().Has("versioning") {
		return GetBucketVersioning(a, w, r)
	}

	if stat.IsDir() {
		return ListObjectsV2(a, w, r)
	}

	if req.URL.Query().Has("versions") {
		return ListObjectVersions(a, w, r)
	}

	return GetObject(a, w, r)
}

func Put(a *S.App, w http.ResponseWriter, req *http.Request) error {
	log.Printf(">>> PUT %v\n", req)

	r, err := a.ParseRequest(req)
	if err != nil {
		return a.RespondError(w, 500, "InternalError", err, r.Bucket)
	}

	if len(r.Key) > 0 {
		if len(req.Header.Get("X-Amz-Copy-Source")) > 0 {
			return CopyObject(a, w, r, req)
		}
		return PutObject(a, w, r, req)
	}

	return CreateBucket(a, w, r)
}

func Post(a *S.App, w http.ResponseWriter, req *http.Request) error {
	log.Printf(">>> POST %v\n", req)

	r, err := a.ParseRequest(req)
	if err != nil {
		return a.RespondError(w, 500, "InternalError", err, "")
	}

	if req.URL.Query().Has("uploads") {
		return CreateMultipartUpload(a, w, r)
	}

	if req.URL.Query().Has("delete") {
		return DeleteObjects(a, w, r, req)
	}

	return a.RespondError(w, 500, "InternalError", err, "")
}

func Delete(a *S.App, w http.ResponseWriter, req *http.Request) error {
	log.Printf(">>> DELETE %v\n", req)

	r, err := a.ParseRequest(req)
	if err != nil {
		return a.RespondError(w, 500, "InternalError", err, r.Bucket)
	}

	if len(r.Key) > 0 {
		return DeleteObject(a, w, r)
	}

	return DeleteBucket(a, w, r)
}

func Head(a *S.App, w http.ResponseWriter, req *http.Request) error {
	log.Printf(">>> HEAD %v\n", req)

	r, err := a.ParseRequest(req)
	if err != nil {
		return a.RespondError(w, 500, "InternalError", err, r.Bucket)
	}

	if len(r.Key) > 0 {
		return HeadObject(a, w, r)
	}
	return HeadBucket(a, w, r)
}
