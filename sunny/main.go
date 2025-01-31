package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String("http://localhost:8080/"),
		Region:           aws.String("us-east-2"), // Required but not used for S3Proxy
		DisableSSL:       aws.Bool(true),          // S3Proxy usually runs on HTTP
		S3ForcePathStyle: aws.Bool(true),          // Equivalent to withPathStyleAccess(true)
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// Create an S3 client
	s3Client := s3.New(sess)

	// Example: List buckets
	result, err := s3Client.ListBuckets(nil)
	if err != nil {
		log.Fatalf("Failed to list buckets: %v", err)
	}

	// Print bucket names
	fmt.Println("Buckets:")
	for _, b := range result.Buckets {
		fmt.Printf("* %s\n", aws.StringValue(b.Name))
	}

	// Define the bucket and object key
	bucketName := "kubestash"
	objectKey := "demo-object.txt"
	objectContent := "Hello, KubeStash! This is a test object."

	// Upload object
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader([]byte(objectContent)),
	})
	if err != nil {
		log.Fatalf("Failed to upload object: %v", err)
	}

	fmt.Printf("Successfully uploaded object '%s' to bucket '%s'\n", objectKey, bucketName)

}
