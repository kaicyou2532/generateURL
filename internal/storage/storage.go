package storage

import (
	"context"
	"mime/multipart"
)

// Saver persists uploaded files and returns a stable identifier that can be used in URLs.
type Saver interface {
	Save(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
}
