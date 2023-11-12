// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"log"
	"net/http"
	"os"

	S "github.com/autovia/s3-go/structs"
)

func ListBuckets(app *S.App, w http.ResponseWriter, req *http.Request) error {
	log.Printf("#ListBuckets %v\n", req)

	files, err := os.ReadDir(*app.Mount)
	if err != nil {
		return app.RespondError(w, 500, "InternalError", "InternalError", "")
	}

	buckets := []S.Bucket{}
	for _, file := range files {
		fileInfo, _ := file.Info()
		if file.IsDir() {
			buckets = append(buckets, S.Bucket{Name: fileInfo.Name(), CreationDate: fileInfo.ModTime()})
		}
	}

	bucketList := S.ListAllMyBucketsResult{
		Owner:   &S.Owner{ID: "123", DisplayName: "jan"},
		Buckets: buckets,
	}

	return app.RespondXML(w, http.StatusOK, bucketList)
}

func CreateBucket(app *S.App, w http.ResponseWriter, r *S.Request) error {
	log.Printf("#CreateBucket: %v\n", r)

	if _, err := os.Stat(r.Path); !os.IsNotExist(err) {
		return app.RespondError(w, 409, "BucketAlreadyExists", "BucketAlreadyExists", r.Bucket)
	}

	if err := os.Mkdir(r.Path, os.ModePerm); err != nil {
		return app.RespondError(w, 500, "InternalError", "InternalError", r.Bucket)
	}

	return app.RespondXML(w, http.StatusOK, nil)
}

func HeadBucket(app *S.App, w http.ResponseWriter, r *S.Request) error {
	log.Printf("#HeadBucket: %v\n", r)

	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return app.RespondError(w, 400, "NoSuchBucket", "NoSuchBucket", r.Bucket)
	}

	return app.Respond(w, http.StatusOK, nil, nil)
}

func GetBucketVersioning(app *S.App, w http.ResponseWriter, r *S.Request) error {
	log.Printf("#GetBucketVersioning: %v\n", r)

	return app.RespondXML(w, http.StatusOK, S.VersioningConfiguration{Status: "Suspended"})
}

func DeleteBucket(app *S.App, w http.ResponseWriter, r *S.Request) error {
	log.Printf("#DeleteBucket: %v\n", r)

	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Bucket)
	}

	contents, err := os.ReadDir(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Bucket)
	}

	if len(contents) > 0 {
		return app.RespondError(w, http.StatusConflict, "BucketNotEmpty", "BucketNotEmpty", r.Bucket)
	}

	if err := os.Remove(r.Path); err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Bucket)
	}

	return app.RespondXML(w, http.StatusNoContent, nil)
}
