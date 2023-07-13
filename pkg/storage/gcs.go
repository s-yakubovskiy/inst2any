package storage

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCS struct {
	bucketName string
	client     *storage.Client
	publicURL  string
}

func NewGCS(bucketName, credsFilePath string) (*GCS, error) {
	ctx := context.Background()

	// Initialize the GCS client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credsFilePath))
	if err != nil {
		return nil, err
	}

	return &GCS{
		bucketName: bucketName,
		client:     client,
		publicURL:  "https://storage.googleapis.com",
	}, nil
}

func (g *GCS) Upload(ctx context.Context, directory, objectName string, r io.Reader) error {
	// Get a handle to the bucket
	bucket := g.client.Bucket(g.bucketName)

	// Get a handle to the object in the bucket
	obj := bucket.Object(fmt.Sprintf("%s/%s", directory, objectName))

	// Prepare the object writer
	w := obj.NewWriter(ctx)

	// Copy the data to the bucket
	if _, err := io.Copy(w, r); err != nil {
		return err
	}

	// Close the writer and check for errors
	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (g *GCS) ReturnPublicURL(ctx context.Context, directory, objectName string) string {
	urlDL := fmt.Sprintf("%s/%s/%s/%s", g.publicURL, g.bucketName, directory, objectName)

	return urlDL
}
