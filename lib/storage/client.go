package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// NewS3Client returns configured AWS S3 client
func NewS3Client(opt Options) *s3.S3 {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(opt.Key, opt.Secret, ""),
		Endpoint:         aws.String(opt.Endpoint),
		Region:           aws.String(opt.Region),
		DisableSSL:       aws.Bool(opt.DisableSSL),
		S3ForcePathStyle: aws.Bool(opt.ForcePathStyle),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		println(err)
	}
	return s3.New(newSession)
}
