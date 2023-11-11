package handlers

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	S "github.com/autovia/s3-go/structs"
)

func Get(a *S.App, w http.ResponseWriter, req *http.Request) error {
	log.Printf(">>> GET %v\n", req)

	if req.URL.Path == "/" {
		return ListBuckets(a, w, req)
	}

	r, err := a.ParseRequest(req)
	if err != nil {
		return a.RespondError(w, 500, "InternalError", "InternalError", r.Bucket)
	}

	stat, err := os.Stat(r.Path)
	if os.IsNotExist(err) {
		return a.RespondError(w, http.StatusInternalServerError, "InternalError", "InternalError", r.Bucket)
	}

	if len(req.URL.Query().Get("versioning")) > 0 {
		return GetBucketVersioning(a, w, r)
	}

	if stat.IsDir() {
		return ListObjectsV2(a, w, r)
	}

	if len(req.URL.Query().Get("versions")) > 0 {
		return ListObjectVersions(a, w, r)
	}

	return GetObject(a, w, r)
}

func Put(a *S.App, w http.ResponseWriter, req *http.Request) error {
	log.Printf(">>> PUT %v\n", req)

	r, err := a.ParseRequest(req)
	if err != nil {
		return a.RespondError(w, 500, "InternalError", "InternalError", r.Bucket)
	}

	if len(r.Key) > 0 {
		source := req.Header.Get("X-Amz-Copy-Source")
		sourcePath, err := url.QueryUnescape(source)
		if err != nil {
			return a.RespondError(w, 500, "InternalError", "InternalError", r.Bucket)
		}
		if len(sourcePath) > 0 {
			path, ok := strings.CutPrefix(sourcePath, "/")
			if !ok {
				return a.RespondError(w, 500, "InternalError", "InternalError", r.Bucket)
			}
			r.Key = sourcePath
			r.Path = strings.Join([]string{*a.Mount, path}, "/")
			return CopyObject(a, w, r, req)
		}
		return PutObject(a, w, r, req)
	}
	return CreateBucket(a, w, r)
}

func Delete(a *S.App, w http.ResponseWriter, req *http.Request) error {
	log.Printf(">>> DELETE %v\n", req)

	r, err := a.ParseRequest(req)
	if err != nil {
		return a.RespondError(w, 500, "InternalError", "InternalError", r.Bucket)
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
		return a.RespondError(w, 500, "InternalError", "InternalError", r.Bucket)
	}

	if len(r.Key) > 0 {
		return HeadObject(a, w, r)
	}
	return HeadBucket(a, w, r)
}
