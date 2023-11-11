// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	S "github.com/autovia/s3-go/structs"
)

const iso8601TimeFormat = "2006-01-02T15:04:05.000Z"

func ListObjectsV2(app *S.App, w http.ResponseWriter, r *S.Request) error {
	log.Printf("#ListObjectsV2 %v\n", r)

	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Bucket)
	}

	contents, err := os.ReadDir(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Bucket)
	}

	objects := []S.Object{}
	prefixes := []S.CommonPrefix{}
	for _, file := range contents {
		fileInfo, _ := file.Info()
		if !file.IsDir() {
			t := fileInfo.ModTime()
			objects = append(objects, S.Object{
				Key:          fileInfo.Name(),
				LastModified: t.Format(iso8601TimeFormat),
				Size:         fileInfo.Size(),
				ETag:         fileInfo.Name(),
				StorageClass: "STANDARD"})
		} else {
			prefixes = append(prefixes, S.CommonPrefix{Prefix: fileInfo.Name() + "/"})
		}
	}

	listBucketResult := S.ListBucketResult{
		Name:           r.Bucket,
		KeyCount:       len(objects),
		MaxKeys:        1000,
		IsTruncated:    false,
		Contents:       objects,
		CommonPrefixes: prefixes,
	}

	return app.RespondXML(w, http.StatusOK, listBucketResult)
}

func CopyObject(app *S.App, w http.ResponseWriter, r *S.Request, req *http.Request) error {
	log.Printf("#CopyObject: %v\n", r)

	sourceFile, err := os.Open(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
	}
	defer targetFile.Close()
	io.Copy(targetFile, sourceFile)

	stats, err := os.Stat(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
	}
	t := stats.ModTime()
	return app.RespondXML(w, http.StatusOK, S.CopyObjectResponse{
		LastModified: t.Format(iso8601TimeFormat),
		ETag:         "123",
	})
}

func PutObject(app *S.App, w http.ResponseWriter, r *S.Request, req *http.Request) error {
	log.Printf("#PutObject: %v\n", r)

	if _, err := os.Stat(r.Path); !os.IsNotExist(err) {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
	}

	if strings.HasSuffix(r.Path, "/") {
		err := os.Mkdir(r.Path, os.ModePerm)
		if err != nil {
			return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
		}
		return app.Respond(w, http.StatusOK, nil, nil)
	}

	targetFile, err := os.Create(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
	}
	defer targetFile.Close()

	defer req.Body.Close()
	io.Copy(targetFile, req.Body)

	return app.Respond(w, http.StatusOK, nil, nil)
}

func HeadObject(app *S.App, w http.ResponseWriter, r *S.Request) error {
	log.Printf("#HeadObject: %v\n", r)

	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return app.RespondError(w, http.StatusBadRequest, "NoSuchKey", "NoSuchKey", r.Key)
	}

	file, err := os.Stat(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
	}

	headers := make(map[string]string)
	t := file.ModTime()
	headers["Content-Length"] = fmt.Sprintf("%v", file.Size())
	headers["Last-Modified"] = t.Format(iso8601TimeFormat)
	headers["ETag"] = "xxx"
	headers["X-Amz-Meta-Autovia"] = "ARCHIVE_ACCESS"

	return app.Respond(w, http.StatusOK, headers, nil)
}

func GetObject(app *S.App, w http.ResponseWriter, r *S.Request) error {
	log.Printf("#GetObject: %v\n", r)

	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return app.RespondError(w, 400, "NoSuchKey", "NoSuchKey", r.Key)
	}

	file, err := os.Open(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
	}
	stats, err := file.Stat()
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
	}
	if stats.IsDir() {
		return app.RespondError(w, 400, "NoSuchKey", "NoSuchKey", r.Key)
	}

	headers := make(map[string]string)
	t := stats.ModTime()
	headers["Content-Length"] = fmt.Sprintf("%v", stats.Size())
	headers["Last-Modified"] = t.Format(iso8601TimeFormat)

	return app.RespondFile(w, http.StatusOK, headers, file)
}

func ListObjectVersions(app *S.App, w http.ResponseWriter, r *S.Request) error {
	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return app.RespondError(w, 400, "NoSuchKey", "NoSuchKey", r.Key)
	}

	file, err := os.Open(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
	}
	stats, err := file.Stat()
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
	}
	t := stats.ModTime()
	return app.RespondXML(w, http.StatusOK, S.ListVersionsResult{
		Name:        r.Bucket,
		Prefix:      r.Key,
		MaxKeys:     1,
		IsTruncated: false,
		Version: []S.ObjectVersion{
			{
				Object: S.Object{
					Key:          r.Key,
					LastModified: t.Format(iso8601TimeFormat),
					ETag:         "xxx",
					Size:         stats.Size(),
					StorageClass: "STANDARD",
					Owner:        &S.Owner{ID: "123", DisplayName: "jan"},
				},
				IsLatest:  true,
				VersionID: "xxx",
			},
		},
	})
}

func DeleteObject(app *S.App, w http.ResponseWriter, r *S.Request) error {
	log.Printf("#DeleteObject: %v\n", r)

	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
	}

	if err := os.Remove(r.Path); err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Key)
	}

	return nil
}
