// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package s3

import (
	"encoding/xml"
	"time"
)

type ListBucketResult struct {
	XMLName        xml.Name `xml:"ListBucketResult"`
	Name           string
	Prefix         string
	KeyCount       int
	MaxKeys        int
	Delimiter      string `xml:"Delimiter,omitempty"`
	IsTruncated    bool
	Contents       []Object
	CommonPrefixes []CommonPrefix
	EncodingType   string `xml:"EncodingType,omitempty"`
}

type CommonPrefix struct {
	Prefix string
}

type Object struct {
	Key          string
	LastModified string
	ETag         string
	Size         int64
	Owner        *Owner `xml:"Owner,omitempty"`
	StorageClass string
}

type Metadata struct {
	Items []struct {
		Key   string
		Value string
	}
}

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Xmlns   string   `xml:"xmlns,attr"`
	Buckets []Bucket `xml:"Buckets>Bucket"`
	Owner   *Owner   `xml:"Owner,omitempty"`
}

type Bucket struct {
	Name         string    `xml:"Name"`
	CreationDate time.Time `xml:"CreationDate"`
}

type Owner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}
