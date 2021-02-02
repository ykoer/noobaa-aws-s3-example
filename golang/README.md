# NooBaa Object Bucket Golang Example
Simple Golang application which shows how you can interact with a NooBaa buckets.

## Build

Golang Binary:

```
go build -a -o ./noobaa-bucket
```

Image:

```
docker build -t noobaa-bucket-example .
```

## Usage

```
./noobaa-bucket-example <command> [<args>

commands:
 list
 upload <key> <source-file-path>
 download <key> <target-file-path>

```

Examples:

```
./noobaa-bucket-example list
./noobaa-bucket-example upload testfile /tmp/testfile
./noobaa-bucket-example download testfile /tmp/testfile-copy
```

Container:

```
docker run -it -e AWS_ACCESS_KEY_ID=******** -e AWS_SECRET_ACCESS_KEY=******** -e BUCKET_HOST=******** -e BUCKET_NAME=******** noobaa-bucket-example:latest list
```

## Required Environment varibles
The environment variables listed below are required in my test application. NooBaa creates a ConfigMap and a Secret with those values in the namespace where the ObjectBucketClaim resource is created.

* AWS\_ACCESS\_KEY\_ID	
* AWS\_SECRET\_ACCESS\_KEY
* BUCKET\_HOST
* BUCKET\_NAME
* BUCKET\_REGION:

