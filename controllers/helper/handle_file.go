package helper

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// Firebase Storage bucket name
const bucketName = "upload-images-f97b5.appspot.com"

// UploadFile uploads a file to Firebase Storage
func UploadFileToCloud(filePath, fileName string) (string, error) {
	// Load Firebase credentials
	ctx := context.Background()
	opt := option.WithCredentialsFile("upload-config.json")

	// Initialize Firebase App
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return "", fmt.Errorf("error initializing Firebase app: %v", err)
	}

	// Initialize Cloud Storage client
	client, err := app.Storage(ctx)
	if err != nil {
		return "", fmt.Errorf("error initializing Cloud Storage: %v", err)
	}

	// Get the storage bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to get bucket: %v", err)
	}

	// Open the file to upload
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Upload the file
	object := bucket.Object(fileName)
	writer := object.NewWriter(ctx)
	if _, err := writer.Write([]byte(filePath)); err != nil {
		return "", fmt.Errorf("failed to write to bucket: %v", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Tạo URL truy cập file
	fileURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName)
	return fileURL, nil
}
