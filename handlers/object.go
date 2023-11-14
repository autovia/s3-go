// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	S "github.com/autovia/s3-go/structs"
)

const ISO8601UTCFormat = "2006-01-02T15:04:05.000Z"
const RFC822Format = "Mon, 2 Jan 2006 15:04:05 GMT"

func ListObjectsV2(app *S.App, w http.ResponseWriter, r *S.Request) error {
	log.Printf("#ListObjectsV2 %v\n", r)

	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Bucket)
	}

	contents, err := os.ReadDir(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Bucket)
	}

	objects := []S.Object{}
	prefixes := []S.CommonPrefix{}
	for _, file := range contents {
		fileInfo, _ := file.Info()
		if !file.IsDir() {
			t := fileInfo.ModTime()
			objects = append(objects, S.Object{
				Key:          fileInfo.Name(),
				LastModified: t.Format(RFC822Format),
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

	source := req.Header.Get("X-Amz-Copy-Source")
	sourcePath, err := url.QueryUnescape(source)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Bucket)
	}

	sourceFile, err := os.Open(*app.Mount + sourcePath)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}
	defer sourceFile.Close()

	if _, err := os.Stat(filepath.Dir(r.Path)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(r.Path), os.ModePerm)
		if err != nil {
			return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
		}
	}

	targetFile, err := os.Create(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}
	defer targetFile.Close()
	io.Copy(targetFile, sourceFile)

	stats, err := os.Stat(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}
	t := stats.ModTime()
	return app.RespondXML(w, http.StatusOK, S.CopyObjectResponse{
		LastModified: t.Format(ISO8601UTCFormat),
		ETag:         "123",
	})
}

func CreateMultipartUpload(app *S.App, w http.ResponseWriter, r *S.Request) error {
	log.Printf("#CreateMultipartUpload: %v\n", r)

	if _, err := os.Stat(r.Path); !os.IsNotExist(err) {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}

	if strings.HasSuffix(r.Path, "/") {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", errors.New("path is a directory"), r.Key)
	}

	uploadID := generate(50)
	metapath := filepath.Join(*app.Mount, *app.Metadata, uploadID)
	if err := os.MkdirAll(metapath, os.ModePerm); err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}

	newfile := filepath.Join(metapath, r.Key)
	if _, err := os.Stat(filepath.Dir(newfile)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(newfile), os.ModePerm)
		if err != nil {
			return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
		}
	}

	f, err := os.Create(newfile)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}
	defer f.Close()

	return app.RespondXML(w, http.StatusOK, S.InitiateMultipartUploadResponse{
		Bucket:   r.Bucket,
		Key:      r.Key,
		UploadID: uploadID,
	})
}

func PutObject(app *S.App, w http.ResponseWriter, r *S.Request, req *http.Request) error {
	log.Printf("#PutObject: %v\n", r)

	if _, err := os.Stat(r.Path); !os.IsNotExist(err) {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}

	if strings.HasSuffix(r.Path, "/") {
		err := os.MkdirAll(r.Path, os.ModePerm)
		if err != nil {
			return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
		}
		return app.Respond(w, http.StatusOK, nil, nil)
	}

	if _, err := os.Stat(filepath.Dir(r.Path)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(r.Path), os.ModePerm)
		if err != nil {
			return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
		}
	}

	targetFile, err := os.Create(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}
	defer targetFile.Close()

	defer req.Body.Close()
	io.Copy(targetFile, req.Body)

	return app.Respond(w, http.StatusOK, nil, nil)
}

func HeadObject(app *S.App, w http.ResponseWriter, r *S.Request) error {
	log.Printf("#HeadObject: %v\n", r)

	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return app.RespondError(w, http.StatusBadRequest, "NoSuchKey", err, r.Key)
	}

	file, err := os.Stat(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}

	headers := make(map[string]string)
	t := file.ModTime()
	headers["Content-Length"] = fmt.Sprintf("%v", file.Size())
	headers["Last-Modified"] = t.Format(RFC822Format)
	headers["ETag"] = "xxx"
	headers["X-Amz-Meta-Autovia"] = "ARCHIVE_ACCESS"

	return app.Respond(w, http.StatusOK, headers, nil)
}

func GetObject(app *S.App, w http.ResponseWriter, r *S.Request) error {
	log.Printf("#GetObject: %v\n", r)

	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return app.RespondError(w, 400, "NoSuchKey", err, r.Key)
	}

	file, err := os.Open(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}
	stats, err := file.Stat()
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}
	if stats.IsDir() {
		return app.RespondError(w, 400, "NoSuchKey", err, r.Key)
	}

	headers := make(map[string]string)
	t := stats.ModTime()
	headers["Content-Length"] = fmt.Sprintf("%v", stats.Size())
	headers["Last-Modified"] = t.Format(RFC822Format)

	return app.RespondFile(w, http.StatusOK, headers, file)
}

func ListObjectVersions(app *S.App, w http.ResponseWriter, r *S.Request) error {
	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return app.RespondError(w, 400, "NoSuchKey", err, r.Key)
	}

	file, err := os.Open(r.Path)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}
	stats, err := file.Stat()
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
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
					LastModified: t.Format(ISO8601UTCFormat),
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
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}

	if err := os.RemoveAll(r.Path); err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}
	headers := make(map[string]string)
	headers["Content-Length"] = "0"

	return app.Respond(w, http.StatusOK, headers, nil)
}

func DeleteObjects(app *S.App, w http.ResponseWriter, r *S.Request, req *http.Request) error {
	log.Printf("#DeleteObjects: %v\n", r)

	body, _ := io.ReadAll(req.Body)
	var delete S.Delete
	err := xml.Unmarshal(body, &delete)
	if err != nil {
		return app.RespondError(w, http.StatusInternalServerError, "InternalError", err, r.Key)
	}

	objects := []S.DeletedObject{}
	errors := []S.DeleteError{}
	for _, file := range delete.Objects {
		path := filepath.Join(*app.Mount, r.Bucket, file.Key)
		delErr := S.DeleteError{}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			delErr = S.DeleteError{
				Code:    "NoSuchKey",
				Message: "NoSuchKey",
				Key:     file.Key,
			}
		}

		if err := os.RemoveAll(path); err != nil {
			delErr = S.DeleteError{
				Code:    "NoSuchKey",
				Message: "NoSuchKey",
				Key:     file.Key,
			}
		}

		if delErr != (S.DeleteError{}) {
			errors = append(errors, delErr)
		} else {
			obj := S.DeletedObject{
				Key: file.Key,
			}
			objects = append(objects, obj)
		}
	}

	log.Print(">>> DEL: ", objects)

	return app.RespondXML(w, http.StatusOK, S.DeleteObjectsResponse{
		DeletedObjects: objects,
		Errors:         errors,
	})
}
