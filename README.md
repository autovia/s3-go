# s3-go

**IMPORTANT:** This is not production-ready software. This project is in active development.

## Introduction

A lightweight S3 server without dependencies.

Supports Authorization Header (AWS Signature Version 4)

Tested with:
* aws-cli/2.13.30 or greater
* aws-sdk-go-v2 v1.22.1
* aws-sdk-ruby3/3.185.2

Compatibility tests with [s3-testsuite](https://github.com/autovia/s3-testsuite)

## Development setup

Configure aws cli

cat $HOME/.aws/config

```shell
[default]
endpoint_url=http://localhost:3000
```

cat $HOME/.aws/config

```shell
[default]
aws_access_key_id = user
aws_secret_access_key = password
region = us-east-1
```

Run server

```shell
go run main.go
```

List buckets

```shell
aws s3 ls

2023-11-02 17:15:01 test-bucket
2023-11-02 17:15:01 test-bucket2
```

## License

[Apache License 2.0](https://github.com/autovia/flightdeck/blob/master/LICENSE)

----
_Copyright [Autovia GmbH](https://autovia.io)_