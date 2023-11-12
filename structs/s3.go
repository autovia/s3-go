// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

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

type VersioningConfiguration struct {
	XMLName xml.Name `xml:"VersioningConfiguration"`
	Status  string   `xml:"Status"`
}

type CopyObjectResponse struct {
	XMLName      xml.Name `xml:"CopyObjectResponse"`
	LastModified string
	ETag         string
}

type ListVersionsResult struct {
	XMLName             xml.Name `xml:"ListVersionsResult"`
	Name                string
	Prefix              string
	KeyMarker           string
	NextKeyMarker       string `xml:"NextKeyMarker,omitempty"`
	NextVersionIDMarker string `xml:"NextVersionIdMarker"`
	VersionIDMarker     string `xml:"VersionIdMarker"`
	MaxKeys             int
	Delimiter           string `xml:"Delimiter,omitempty"`
	IsTruncated         bool
	CommonPrefixes      []CommonPrefix
	Version             []ObjectVersion
	EncodingType        string `xml:"EncodingType,omitempty"`
}

type ObjectVersion struct {
	Object
	IsLatest       bool
	VersionID      string `xml:"VersionId"`
	isDeleteMarker bool
}

type Delete struct {
	Objects []Object `xml:"Object"`
	Quiet   bool
}

type DeleteObjectsResponse struct {
	XMLName        xml.Name        `xml:"DeleteObjectsResponse"`
	DeletedObjects []DeletedObject `xml:"Deleted,omitempty"`
	Errors         []DeleteError   `xml:"Error,omitempty"`
}

type DeleteError struct {
	Code      string
	Message   string
	Key       string
	VersionID string `xml:"VersionId"`
}

type DeletedObject struct {
	DeleteMarker          bool   `xml:"DeleteMarker,omitempty"`
	DeleteMarkerVersionID string `xml:"DeleteMarkerVersionId,omitempty"`
	Key                   string `xml:"Key,omitempty"`
	VersionID             string `xml:"VersionId,omitempty"`
}
