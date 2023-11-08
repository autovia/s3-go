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

	S "github.com/autovia/s3/structs"
	"github.com/autovia/s3/structs/s3"
)

const iso8601TimeFormat = "2006-01-02T15:04:05.000Z"

func ListObjectsV2(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#ListObjectsV2 %v\n", r)

	query := r.URL.Query()
	bucket := r.PathValue("bucket")
	path := fmt.Sprintf("%s/%s", *app.Mount, bucket)
	if len(query["prefix"]) > 0 {
		path += "/" + query["prefix"][0]
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s3.RespondError(w, 500, "InternalError", "InternalError", bucket)
	}

	contents, err := os.ReadDir(path)
	if err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", bucket)
	}

	objects := []s3.Object{}
	prefixes := []s3.CommonPrefix{}
	for _, file := range contents {
		fileInfo, _ := file.Info()
		if !file.IsDir() {
			t := fileInfo.ModTime()
			objects = append(objects, s3.Object{
				Key:          fileInfo.Name(),
				LastModified: t.Format(iso8601TimeFormat),
				Size:         fileInfo.Size(),
				ETag:         fileInfo.Name(),
				StorageClass: "STANDARD"})
		} else {
			prefixes = append(prefixes, s3.CommonPrefix{Prefix: fileInfo.Name() + "/"})
		}
	}

	listBucketResult := s3.ListBucketResult{
		Name:           bucket,
		KeyCount:       len(objects),
		MaxKeys:        1000,
		IsTruncated:    false,
		Contents:       objects,
		CommonPrefixes: prefixes,
	}

	return app.RespondXML(w, r, listBucketResult, http.StatusOK)
}

func PutObject(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#PutObject: %v\n", r)

	name, found := strings.CutPrefix(r.URL.Path, "/")
	if !found {
		return nil
	}

	path := fmt.Sprintf("%s/%s", *app.Mount, name)
	log.Print(path)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	defer r.Body.Close()
	out, err := os.Create(path)
	if err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}
	defer out.Close()
	io.Copy(out, r.Body)

	return app.Respond(w, 200, nil, nil)
}

func HeadObject(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#HeadObject: %v\n", r)

	name, found := strings.CutPrefix(r.URL.Path, "/")
	if !found {
		return nil
	}

	path := fmt.Sprintf("%s/%s", *app.Mount, name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s3.RespondError(w, 400, "NoSuchKey", "NoSuchKey", name)
	}

	file, err := os.Stat(path)
	if err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	headers := make(map[string]string)
	t := file.ModTime()
	headers["Content-Length"] = fmt.Sprintf("%v", file.Size())
	headers["Last-Modified"] = t.Format(iso8601TimeFormat)

	return app.Respond(w, 200, headers, nil)
}

func GetObject(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#GetObject: %v\n", r)

	name, found := strings.CutPrefix(r.URL.Path, "/")
	if !found {
		return nil
	}

	path := fmt.Sprintf("%s/%s", *app.Mount, name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s3.RespondError(w, 400, "NoSuchKey", "NoSuchKey", name)
	}

	file, err := os.Open(path)
	if err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}
	stats, err := file.Stat()
	if err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	headers := make(map[string]string)
	t := stats.ModTime()
	headers["Content-Length"] = fmt.Sprintf("%v", stats.Size())
	headers["Last-Modified"] = t.Format(iso8601TimeFormat)

	return app.RespondFile(w, 200, headers, file)
}

func DeleteObject(app *S.App, w http.ResponseWriter, r *http.Request) error {
	log.Printf("#DeleteObject: %v\n", r)

	name, found := strings.CutPrefix(r.URL.Path, "/")
	if !found {
		return nil
	}

	path := fmt.Sprintf("%s/%s", *app.Mount, name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	if err := os.Remove(path); err != nil {
		return s3.RespondError(w, 500, "InternalError", "InternalError", name)
	}

	return nil
}
