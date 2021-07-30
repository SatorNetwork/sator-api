package storage

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
)

func TestACL_String(t *testing.T) {
	tests := []struct {
		name string
		a    ACL
		want string
	}{
		{"success", Public, "public-read"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("ACL.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	s3mock := &S3Mock{}
	opt := Options{
		Bucket:         "test",
		ForcePathStyle: false,
	}
	type args struct {
		c   Client
		opt Options
	}
	tests := []struct {
		name string
		args args
		want *Interactor
	}{
		{"success", args{s3mock, opt}, &Interactor{client: s3mock, bucket: opt.Bucket}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.c, tt.args.opt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInteractor_Upload(t *testing.T) {
	s3mock := &S3Mock{}
	s3mockErr := &S3Mock{Error: errors.New("upload")}
	opt := Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	type fields struct {
		client         Client
		bucket         string
		url            string
		forcePathStyle bool
	}
	type args struct {
		file     io.ReadSeeker
		filepath string
		acl      ACL
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"success", fields{s3mock, opt.Bucket, opt.URL, opt.ForcePathStyle}, args{nil, "test.png", Public}, false},
		{"error", fields{s3mockErr, opt.Bucket, opt.URL, opt.ForcePathStyle}, args{nil, "test.png", Public}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				client:         tt.fields.client,
				bucket:         tt.fields.bucket,
				url:            tt.fields.url,
				forcePathStyle: tt.fields.forcePathStyle,
			}
			if err := i.Upload(tt.args.file, tt.args.filepath, tt.args.acl, "image/png"); (err != nil) != tt.wantErr {
				t.Errorf("Interactor.Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInteractor_Download(t *testing.T) {
	s3mock := &S3Mock{}
	s3mockErr := &S3Mock{Error: errors.New("upload")}
	opt := Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	obj := &s3.GetObjectOutput{}
	type fields struct {
		client         Client
		bucket         string
		url            string
		forcePathStyle bool
	}
	type args struct {
		filepath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.ReadCloser
		wantErr bool
	}{
		{"success", fields{s3mock, opt.Bucket, opt.URL, opt.ForcePathStyle}, args{"test.png"}, obj.Body, false},
		{"error", fields{s3mockErr, opt.Bucket, opt.URL, opt.ForcePathStyle}, args{"test.png"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				client:         tt.fields.client,
				bucket:         tt.fields.bucket,
				url:            tt.fields.url,
				forcePathStyle: tt.fields.forcePathStyle,
			}
			got, _, err := i.Download(tt.args.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interactor.Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Interactor.Download() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInteractor_Remove(t *testing.T) {
	s3mock := &S3Mock{}
	s3mockErr := &S3Mock{Error: errors.New("upload")}
	opt := Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}

	type fields struct {
		client         Client
		bucket         string
		url            string
		forcePathStyle bool
	}
	type args struct {
		filepath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"success", fields{s3mock, opt.Bucket, opt.URL, opt.ForcePathStyle}, args{"test.png"}, false},
		{"error", fields{s3mockErr, opt.Bucket, opt.URL, opt.ForcePathStyle}, args{"test.png"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				client:         tt.fields.client,
				bucket:         tt.fields.bucket,
				url:            tt.fields.url,
				forcePathStyle: tt.fields.forcePathStyle,
			}
			if err := i.Remove(tt.args.filepath); (err != nil) != tt.wantErr {
				t.Errorf("Interactor.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInteractor_FilePath(t *testing.T) {
	type fields struct {
		client         *s3.S3
		bucket         string
		url            string
		forcePathStyle bool
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"success", fields{nil, "test", "http://test", true}, args{"file.png"}, "uploads/file.png"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				client:         tt.fields.client,
				bucket:         tt.fields.bucket,
				url:            tt.fields.url,
				forcePathStyle: tt.fields.forcePathStyle,
			}
			if got := i.FilePath(tt.args.filename); got != tt.want {
				t.Errorf("Interactor.FilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInteractor_FileURL(t *testing.T) {
	type fields struct {
		client         *s3.S3
		bucket         string
		url            string
		forcePathStyle bool
	}
	type args struct {
		filepath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"force path style: true", fields{nil, "test", "http://test", true}, args{"file.png"}, "http://test/test/file.png"},
		{"force path style: false", fields{nil, "test", "http://test", false}, args{"file.png"}, "http://test/file.png"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				client:         tt.fields.client,
				bucket:         tt.fields.bucket,
				url:            tt.fields.url,
				forcePathStyle: tt.fields.forcePathStyle,
			}
			if got := i.FileURL(tt.args.filepath); got != tt.want {
				t.Errorf("Interactor.FileURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
