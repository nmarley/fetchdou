package main

import (
	// "bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Put writes bytes to an object in an s3 bucket
// acls: private, public-read
func S3Put(r io.ReadSeeker, bucket, s3key, s3acl string, sess *session.Session) error {
	if _, err := s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3key),
		ACL:    aws.String(s3acl),
		Body:   r,
	}); err != nil {
		return err
	}
	return nil
}

// S3Get gets an object from an s3 bucket and returns a readseeker
func S3Get(bucket, s3key string, sess *session.Session) (io.ReadCloser, error) {
	out, err := s3.New(sess).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3key),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}
