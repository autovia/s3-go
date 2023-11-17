// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"net/url"
	"reflect"
	"testing"
)

func TestValidSignatureV4GetObject(t *testing.T) {
	headerStr := "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request, SignedHeaders=host;range;x-amz-content-sha256;x-amz-date, Signature=f0e8bdb87c964420e857bd35b5d6ed310bd44f0170aba48dd91039c6036bdb41"

	expectedAuthorizationHeader := map[string]string{
		"Credential":           "AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request",
		"SignedHeaders":        "host;range;x-amz-content-sha256;x-amz-date",
		"Signature":            "f0e8bdb87c964420e857bd35b5d6ed310bd44f0170aba48dd91039c6036bdb41",
		"host":                 "examplebucket.s3.amazonaws.com",
		"range":                "bytes=0-9",
		"x-amz-content-sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"x-amz-date":           "20130524T000000Z",
	}

	header := map[string][]string{
		"Range":                {"bytes=0-9"},
		"X-Amz-Content-Sha256": {"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		"X-Amz-Date":           {"20130524T000000Z"},
	}

	t.Run("authorizationHeader", func(t *testing.T) {
		result := authorizationHeader(header, "examplebucket.s3.amazonaws.com", headerStr)
		if !reflect.DeepEqual(expectedAuthorizationHeader, result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedAuthorizationHeader)
		}
	})

	expectedCanonicalRequest := `GET
/test.txt

host:examplebucket.s3.amazonaws.com
range:bytes=0-9
x-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
x-amz-date:20130524T000000Z

host;range;x-amz-content-sha256;x-amz-date
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`

	t.Run("canonicalRequest", func(t *testing.T) {
		headers := authorizationHeader(header, "examplebucket.s3.amazonaws.com", headerStr)
		u, _ := url.Parse("/test.txt")
		result := canonicalRequest("GET", u, "", []byte(""), headers)
		if expectedCanonicalRequest != result {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedCanonicalRequest)
		}
	})

	expectedStringToSign := `AWS4-HMAC-SHA256
20130524T000000Z
20130524/us-east-1/s3/aws4_request
7344ae5b7ee6c3e7e6b0fe0640412a37625d1fbfff95c48bbb2dc43964946972`

	t.Run("stringToSign", func(t *testing.T) {
		result := stringToSign(expectedCanonicalRequest, "AKIAIOSFODNN7EXAMPLE", expectedAuthorizationHeader)
		if expectedStringToSign != result {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedStringToSign)
		}
	})

	expectedSigningKey := "f0e8bdb87c964420e857bd35b5d6ed310bd44f0170aba48dd91039c6036bdb41"

	t.Run("signingKeySignature", func(t *testing.T) {
		result := signingKeySignature("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", expectedStringToSign, expectedAuthorizationHeader)
		if expectedSigningKey != string(result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedSigningKey)
		}
	})
}

func TestValidSignatureV4PutObject(t *testing.T) {
	headerStr := "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request, SignedHeaders=date;host;x-amz-content-sha256;x-amz-date;x-amz-storage-class, Signature=98ad721746da40c64f1a55b78f14c238d841ea1380cd77a1b5971af0ece108bd"

	expectedAuthorizationHeader := map[string]string{
		"Credential":           "AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request",
		"SignedHeaders":        "date;host;x-amz-content-sha256;x-amz-date;x-amz-storage-class",
		"Signature":            "98ad721746da40c64f1a55b78f14c238d841ea1380cd77a1b5971af0ece108bd",
		"date":                 "Fri, 24 May 2013 00:00:00 GMT",
		"host":                 "examplebucket.s3.amazonaws.com",
		"x-amz-content-sha256": "44ce7dd67c959e0d3524ffac1771dfbba87d2b6b4b4e99e42034a8b803f8b072",
		"x-amz-date":           "20130524T000000Z",
		"x-amz-storage-class":  "REDUCED_REDUNDANCY",
	}

	header := map[string][]string{
		"Date":                 {"Fri, 24 May 2013 00:00:00 GMT"},
		"X-Amz-Content-Sha256": {"44ce7dd67c959e0d3524ffac1771dfbba87d2b6b4b4e99e42034a8b803f8b072"},
		"X-Amz-Date":           {"20130524T000000Z"},
		"x-Amz-Storage-Class":  {"REDUCED_REDUNDANCY"},
	}

	t.Run("authorizationHeader", func(t *testing.T) {
		result := authorizationHeader(header, "examplebucket.s3.amazonaws.com", headerStr)
		if !reflect.DeepEqual(expectedAuthorizationHeader, result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedAuthorizationHeader)
		}
	})

	expectedCanonicalRequest := `PUT
/test%24file.text

date:Fri, 24 May 2013 00:00:00 GMT
host:examplebucket.s3.amazonaws.com
x-amz-content-sha256:44ce7dd67c959e0d3524ffac1771dfbba87d2b6b4b4e99e42034a8b803f8b072
x-amz-date:20130524T000000Z
x-amz-storage-class:REDUCED_REDUNDANCY

date;host;x-amz-content-sha256;x-amz-date;x-amz-storage-class
44ce7dd67c959e0d3524ffac1771dfbba87d2b6b4b4e99e42034a8b803f8b072`

	t.Run("canonicalRequest", func(t *testing.T) {
		headers := authorizationHeader(header, "examplebucket.s3.amazonaws.com", headerStr)
		u, _ := url.Parse("/test%24file.text")
		result := canonicalRequest("PUT", u, "", []byte("Welcome to Amazon S3."), headers)
		if expectedCanonicalRequest != result {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedCanonicalRequest)
		}
	})

	expectedStringToSign := `AWS4-HMAC-SHA256
20130524T000000Z
20130524/us-east-1/s3/aws4_request
9e0e90d9c76de8fa5b200d8c849cd5b8dc7a3be3951ddb7f6a76b4158342019d`

	t.Run("stringToSign", func(t *testing.T) {
		result := stringToSign(expectedCanonicalRequest, "AKIAIOSFODNN7EXAMPLE", expectedAuthorizationHeader)
		if expectedStringToSign != result {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedStringToSign)
		}
	})

	expectedSigningKey := "98ad721746da40c64f1a55b78f14c238d841ea1380cd77a1b5971af0ece108bd"

	t.Run("signingKeySignature", func(t *testing.T) {
		result := signingKeySignature("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", expectedStringToSign, expectedAuthorizationHeader)
		if expectedSigningKey != string(result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedSigningKey)
		}
	})
}

func TestValidSignatureV4BucketLifecycle(t *testing.T) {
	headerStr := "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=f0e8bdb87c964420e857bd35b5d6ed310bd44f0170aba48dd91039c6036bdb41"

	expectedAuthorizationHeader := map[string]string{
		"Credential":           "AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request",
		"SignedHeaders":        "host;x-amz-content-sha256;x-amz-date",
		"Signature":            "f0e8bdb87c964420e857bd35b5d6ed310bd44f0170aba48dd91039c6036bdb41",
		"host":                 "examplebucket.s3.amazonaws.com",
		"x-amz-content-sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"x-amz-date":           "20130524T000000Z",
	}

	header := map[string][]string{
		"X-Amz-Content-Sha256": {"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		"X-Amz-Date":           {"20130524T000000Z"},
	}

	t.Run("authorizationHeader", func(t *testing.T) {
		result := authorizationHeader(header, "examplebucket.s3.amazonaws.com", headerStr)
		if !reflect.DeepEqual(expectedAuthorizationHeader, result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedAuthorizationHeader)
		}
	})

	expectedCanonicalRequest := `GET
/
lifecycle=
host:examplebucket.s3.amazonaws.com
x-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
x-amz-date:20130524T000000Z

host;x-amz-content-sha256;x-amz-date
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`

	t.Run("canonicalRequest", func(t *testing.T) {
		headers := authorizationHeader(header, "examplebucket.s3.amazonaws.com", headerStr)
		u, _ := url.Parse("/")
		result := canonicalRequest("GET", u, "lifecycle=", []byte(""), headers)
		if expectedCanonicalRequest != result {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedCanonicalRequest)
		}
	})

	expectedStringToSign := `AWS4-HMAC-SHA256
20130524T000000Z
20130524/us-east-1/s3/aws4_request
9766c798316ff2757b517bc739a67f6213b4ab36dd5da2f94eaebf79c77395ca`

	t.Run("stringToSign", func(t *testing.T) {
		result := stringToSign(expectedCanonicalRequest, "AKIAIOSFODNN7EXAMPLE", expectedAuthorizationHeader)
		if expectedStringToSign != result {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedStringToSign)
		}
	})

	expectedSigningKey := "fea454ca298b7da1c68078a5d1bdbfbbe0d65c699e0f91ac7a200a0136783543"

	t.Run("signingKeySignature", func(t *testing.T) {
		result := signingKeySignature("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", expectedStringToSign, expectedAuthorizationHeader)
		if expectedSigningKey != string(result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedSigningKey)
		}
	})
}

func TestValidSignatureV4BucketListObjects(t *testing.T) {
	headerStr := "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=34b48302e7b5fa45bde8084f4b7868a86f0a534bc59db6670ed5711ef69dc6f7"

	expectedAuthorizationHeader := map[string]string{
		"Credential":           "AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request",
		"SignedHeaders":        "host;x-amz-content-sha256;x-amz-date",
		"Signature":            "34b48302e7b5fa45bde8084f4b7868a86f0a534bc59db6670ed5711ef69dc6f7",
		"host":                 "examplebucket.s3.amazonaws.com",
		"x-amz-content-sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"x-amz-date":           "20130524T000000Z",
	}

	header := map[string][]string{
		"X-Amz-Content-Sha256": {"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		"X-Amz-Date":           {"20130524T000000Z"},
	}

	t.Run("authorizationHeader", func(t *testing.T) {
		result := authorizationHeader(header, "examplebucket.s3.amazonaws.com", headerStr)
		if !reflect.DeepEqual(expectedAuthorizationHeader, result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedAuthorizationHeader)
		}
	})

	expectedCanonicalRequest := `GET
/
max-keys=2&prefix=J
host:examplebucket.s3.amazonaws.com
x-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
x-amz-date:20130524T000000Z

host;x-amz-content-sha256;x-amz-date
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`

	t.Run("canonicalRequest", func(t *testing.T) {
		headers := authorizationHeader(header, "examplebucket.s3.amazonaws.com", headerStr)
		u, _ := url.Parse("/")
		result := canonicalRequest("GET", u, "max-keys=2&prefix=J", []byte(""), headers)
		if expectedCanonicalRequest != result {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedCanonicalRequest)
		}
	})

	expectedStringToSign := `AWS4-HMAC-SHA256
20130524T000000Z
20130524/us-east-1/s3/aws4_request
df57d21db20da04d7fa30298dd4488ba3a2b47ca3a489c74750e0f1e7df1b9b7`

	t.Run("stringToSign", func(t *testing.T) {
		result := stringToSign(expectedCanonicalRequest, "AKIAIOSFODNN7EXAMPLE", expectedAuthorizationHeader)
		if expectedStringToSign != result {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedStringToSign)
		}
	})

	expectedSigningKey := "34b48302e7b5fa45bde8084f4b7868a86f0a534bc59db6670ed5711ef69dc6f7"

	t.Run("signingKeySignature", func(t *testing.T) {
		result := signingKeySignature("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", expectedStringToSign, expectedAuthorizationHeader)
		if expectedSigningKey != string(result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expectedSigningKey)
		}
	})

}
