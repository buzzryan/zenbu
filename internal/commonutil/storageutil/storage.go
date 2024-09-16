package storageutil

import "context"

type Scope int

const (
	// Private scope means that the file is only accessible if the user has the correct permission.
	// Private file not uploaded to public CDN and public storage.
	Private Scope = iota

	// Public scope means that the file is accessible to anyone.
	// Public file uploaded to public CDN or public storage.
	Public
)

type Storage interface {
	// Many storage services provide a way to generate a signed URL for uploading an object.
	// AWS S3 supports this features by using `pre-signed URLs`.
	// Google cloud storage also supports this feature by using `signed URLs`.
	// Azure Blob storage also supports this feature by using `shared access signatures`.

	// CreateUploadURL returns a signed URL for uploading a file to the storage.
	// filepath is the path of the file to be uploaded.
	CreateUploadURL(ctx context.Context, scope Scope, filepath string) (url string, err error)

	// GetPublicFileURL returns a public URL for accessing a file in the storage.
	GetPublicFileURL(ctx context.Context, filepath string) (url string, err error)
}
