// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package s3

import (
	"encoding/xml"
	"time"
)

type ListAllMyBucketsResult struct {
	XMLName xml.Name  `xml:"ListAllMyBucketsResult"`
	Xmlns   string    `xml:"xmlns,attr"`
	Buckets []Bucket  `xml:"Buckets>Bucket"`
	Owner   *UserInfo `xml:"Owner,omitempty"`
}

type Bucket struct {
	Name         string    `xml:"Name"`
	CreationDate time.Time `xml:"CreationDate"`
}

type UserInfo struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

type CreateBucketConfiguration struct {
	XMLName            xml.Name `xml:"CreateBucketConfiguration"`
	Xmlns              string   `xml:"xmlns,attr"`
	LocationConstraint string   `xml:"LocationConstraint"`
}
