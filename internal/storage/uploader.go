package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// UploadToGCS uploads a file to a GCS bucket and returns the GCS URI.
func UploadToGCS(ctx context.Context, localFilePath, bucketName, credentialsFile string) (string, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return "", fmt.Errorf("storage client oluşturulamadı: %w", err)
	}
	defer client.Close()

	file, err := os.Open(localFilePath)
	if err != nil {
		return "", fmt.Errorf("yerel dosya açılamadı: %w", err)
	}
	defer file.Close()

	// Create a unique object name
	objectName := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(localFilePath))

	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)

	// Upload the file
	wc := object.NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("dosya GCS'ye kopyalanamadı: %w", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("GCS writer kapatılamadı: %w", err)
	}

	gcsURI := fmt.Sprintf("gs://%s/%s", bucketName, objectName)
	return gcsURI, nil
}