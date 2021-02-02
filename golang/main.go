package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	bucketHost   = os.Getenv("BUCKET_HOST")
	bucketName   = os.Getenv("BUCKET_NAME")
	bucketRegion = os.Getenv("BUCKET_REGION")
)

var sess = connectAWS()

func connectAWS() *session.Session {

	s3CustResolverFn := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.S3ServiceID {
			return endpoints.ResolvedEndpoint{
				URL:           bucketHost,
				SigningRegion: "custom-signing-region",
			}, nil
		}

		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	sess, err := session.NewSession(&aws.Config{
		HTTPClient:       client,
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(bucketRegion),
		EndpointResolver: endpoints.ResolverFunc(s3CustResolverFn),
	})

	if err != nil {
		panic(err)
	}
	return sess
}

func main() {

	if len(os.Args) == 1 {
		fmt.Printf("usage: %s <command> [<args>]\n", os.Args[0])
		fmt.Println()
		fmt.Println("commands:")
		fmt.Println(" list")
		fmt.Println(" upload <key> <source-file-path>")
		fmt.Println(" download <key> <target-file-path>")
		return
	}

	switch os.Args[1] {
	case "list":
		err := list()
		if err != nil {
			panic(err)
		}
	case "upload":
		if len(os.Args) == 4 {
			err := upload(os.Args[2], os.Args[3])
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("key and source-file-path are required!")
		}
	case "download":
		if len(os.Args) == 4 {
			err := download(os.Args[2], os.Args[3])
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("key and target-file-path are required!")
		}
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

}

func list() error {

	fmt.Printf("Listing Files in the Bucket: %s\n", bucketName)

	svc := s3.New(sess)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		return fmt.Errorf("failed to list objects: %s", err)
	}
	for _, item := range result.Contents {
		fmt.Printf("File: %s\n", *item.Key)
	}
	return nil
}

func upload(key, sourceFilepath string) error {

	fmt.Printf("Uploading File with the key: %s into the Bucket: %s\n", key, bucketName)

	data, err := ioutil.ReadFile(sourceFilepath)
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	if _, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
	}); err != nil {
		return fmt.Errorf("failed to store object %s: %s", key, err)
	}
	return nil
}

func download(key, targetFilepath string) error {

	fmt.Printf("Downloading File with the key: %s from the Bucket: %s\n", key, bucketName)

	f, err := os.Create(targetFilepath)
	if err != nil {
		return fmt.Errorf("Something went wrong creating the local file: %s", err)
	}

	// Write the contents of S3 Object to the file
	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("Something went wrong retrieving the file from S3: %s", err)
	}

	return nil
}
