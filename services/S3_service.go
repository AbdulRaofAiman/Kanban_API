package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"kanban-backend/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct{}

func NewS3Service() *S3Service {
	return &S3Service{}
}

// UploadFile uploads a file to S3 and returns the public URL
func (s *S3Service) UploadFile(file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	bucket := os.Getenv("AWS_S3_BUCKET")
	region := os.Getenv("AWS_REGION")

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s/%d%s", folder, time.Now().UnixNano(), ext)

	uploader := manager.NewUploader(config.S3Client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return public URL
	_ = result
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, filename)
	return url, nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(key string) error {
	bucket := os.Getenv("AWS_S3_BUCKET")

	_, err := config.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetSignedURL generates a temporary signed URL (expires in 1 hour)
func (s *S3Service) GetSignedURL(key string) (string, error) {
	bucket := os.Getenv("AWS_S3_BUCKET")

	presignClient := s3.NewPresignClient(config.S3Client)

	result, err := presignClient.PresignGetObject(context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(time.Hour),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	return result.URL, nil
}
