// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"log"
	"net/http"
	"os"

	S "github.com/autovia/s3-go/structs"
	"github.com/autovia/s3-go/structs/s3"
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

	return app.RespondXML(w, http.StatusOK, bucketList)
}

func CreateBucket(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#CreateBucket: %v\n", r)

	name, path, err := app.ParseRequest(r)
	if err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return s3.RespondError(w, 409, "BucketAlreadyExists", "BucketAlreadyExists", name)
	}

	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	return app.RespondXML(w, http.StatusOK, nil)
}

func HeadBucket(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#HeadBucket: %v\n", r)

	name, path, err := app.ParseRequest(r)
	if err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s3.RespondError(w, 400, "NoSuchBucket", "NoSuchBucket", name)
	}

	return app.Respond(w, http.StatusOK, nil, nil)
}

func DeleteBucket(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#DeleteBucket: %v\n", r)

	name, path, err := app.ParseRequest(r)
	if err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s3.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", name)
	}

	contents, err := os.ReadDir(path)
	if err != nil {
		return s3.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", name)
	}

	if len(contents) > 0 {
		return s3.RespondError(w, http.StatusConflict, "BucketNotEmpty", "BucketNotEmpty", name)
	}

	if err := os.Remove(path); err != nil {
		return s3.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", name)
	}

	return app.RespondXML(w, http.StatusNoContent, nil)
}
