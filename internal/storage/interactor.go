package storage

import (
	"fmt"
	"io"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

// Predefined ACL permissions
const (
	Public  ACL = "public-read"
	Private ACL = "private"
)

type (
	// ACL permission
	ACL string

	// Interactor struct
	Interactor struct {
		client         Client
		bucket         string
		url            string
		forcePathStyle bool
	}

	// Client interface
	Client interface {
		PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
		GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
		DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	}
)

func (a ACL) String() string {
	return string(a)
}

// New is a factory function,
// returns a new instance of the storage interactor
func New(c Client, opt Options) *Interactor {
	return &Interactor{
		client:         c,
		bucket:         opt.Bucket,
		url:            opt.URL,
		forcePathStyle: opt.ForcePathStyle,
	}
}

// Upload file to the cloud storage
func (i *Interactor) Upload(file io.ReadSeeker, filepath string, acl ACL, contentType string) error {
	object := s3.PutObjectInput{
		Bucket:      aws.String(i.bucket),
		Key:         aws.String(filepath),
		Body:        file,
		ACL:         aws.String(acl.String()),
		ContentType: aws.String(contentType),
	}

	_, err := i.client.PutObject(&object)
	if err != nil {
		return errors.Wrap(err, "storage.upload")
	}

	return nil
}

// Download file from the cloud storage
func (i *Interactor) Download(filepath string) (io.ReadCloser, *string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(filepath),
	}

	result, err := i.client.GetObject(input)
	if err != nil {
		return nil, nil, errors.Wrap(err, "storage.download")
	}

	return result.Body, result.ContentType, nil
}

// Remove file from the cloud storage
func (i *Interactor) Remove(filepath string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(filepath),
	}

	_, err := i.client.DeleteObject(input)
	if err != nil {
		return errors.Wrap(err, "storage.remove")
	}

	return nil
}

// FilePath returns absolute file path in the storage
func (i *Interactor) FilePath(filename string) string {
	return path.Join("uploads", filename)
}

// FileURL return public url for a file
func (i *Interactor) FileURL(filepath string) string {
	if i.forcePathStyle {
		return fmt.Sprintf("%s/%s/%s", i.url, i.bucket, filepath)
	}
	return fmt.Sprintf("%s/%s", i.url, filepath)
}
