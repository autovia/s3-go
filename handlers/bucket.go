// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	S "github.com/autovia/s3/structs"
	"github.com/autovia/s3/structs/s3"
)

func ListBuckets(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#ListBuckets %v\n", r)

	files, err := os.ReadDir(*app.Mount)
	if err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", "")
	}

	buckets := []s3.Bucket{}
	for _, file := range files {
		fileInfo, _ := file.Info()
		if file.IsDir() {
			buckets = append(buckets, s3.Bucket{Name: fileInfo.Name(), CreationDate: fileInfo.ModTime()})
		}
	}

	bucketList := s3.ListAllMyBucketsResult{
		Owner:   &s3.Owner{ID: "123", DisplayName: "jan"},
		Buckets: buckets,
	}

	return app.RespondXML(w, r, bucketList, http.StatusOK)
}

func CreateBucket(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#CreateBucket: %v\n", r)

	name, found := strings.CutPrefix(r.URL.Path, "/")
	if !found {
		return nil
	}

	path := fmt.Sprintf("%s/%s", *app.Mount, name)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return s3.RespondError(w, 409, "BucketAlreadyExists", "BucketAlreadyExists", name)
	}

	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	return app.RespondXML(w, r, nil, http.StatusOK)
}

func HeadBucket(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#HeadBucket: %v\n", r)

	name, found := strings.CutPrefix(r.URL.Path, "/")
	if !found {
		return nil
	}

	path := fmt.Sprintf("%s/%s", *app.Mount, name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s3.RespondError(w, 400, "NoSuchBucket", "NoSuchBucket", name)
	}

	return app.Respond(w, 200, nil, nil)
}

func DeleteBucket(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#DeleteBucket: %v\n", r)

	name, found := strings.CutPrefix(r.URL.Path, "/")
	if !found {
		return nil
	}

	path := fmt.Sprintf("%s/%s", *app.Mount, name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	contents, err := os.ReadDir(path)
	if err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	if len(contents) > 0 {
		return s3.RespondError(w, 409, "BucketNotEmpty", "BucketNotEmpty", name)
	}

	if err := os.Remove(path); err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	return app.RespondXML(w, r, nil, http.StatusNoContent)
}
