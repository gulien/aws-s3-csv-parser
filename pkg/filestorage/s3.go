package filestorage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Downloader implements Downloader for AWS S3.
type S3Downloader struct {
	session *session.Session

	bucket string
	key    string
}

// NewS3PublicDownloader creates an instance of S3Downloader to a public AWS
// S3.
func NewS3PublicDownloader(region, bucket, key string) Downloader {
	return S3Downloader{
		session: session.Must(session.NewSession(&aws.Config{Region: aws.String(region)})),
		bucket:  bucket,
		key:     key,
	}
}

// Download downloads a file from a public AWS S3.
func (storage S3Downloader) Download(ctx context.Context, w io.WriterAt) error {
	downloader := s3manager.NewDownloader(storage.session)

	_, err := downloader.DownloadWithContext(ctx, w, &s3.GetObjectInput{
		Bucket: aws.String(storage.bucket),
		Key:    aws.String(storage.key),
	})

	return err
}

// Interface guards.
var (
	_ Downloader = (*S3Downloader)(nil)
)
