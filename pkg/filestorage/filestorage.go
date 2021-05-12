package filestorage

import (
	"context"
	"io"
)

// Downloader is an interface for downloading files from Cloud storages.
type Downloader interface {
	Download(ctx context.Context, w io.WriterAt) error
}
