package service

import (
	// Standart
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3ServiceConfig struct {
	isLoggingEnabled bool
}

type S3Service struct {
	s3Client *s3.Client
	bucket   string
	config   S3ServiceConfig
}

func NewS3Service() *S3Service {
	return &S3Service{
		s3Client: InitS3Client(),
		bucket:   os.Getenv("S3_BUCKET_NAME"),
		config: S3ServiceConfig{
			isLoggingEnabled: true,
		},
	}
}

func InitS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}

	return s3.NewFromConfig(cfg)
}

// Behaviours --------------------------------------------------------------------

func (s *S3Service) GetObject(path *string, fileName *string, processID string) ([]byte, map[string]string, error) {
	// Define the full key
	key := *path + *fileName

	// util.LogTask("get requesting on S3 with:"+"| filename: "+*fileName+"| path: "+*path+"| key: "+key, "S3Service.GetObject()", "")

	// Set up the GetObject input
	input := &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	}

	// Call S3 GetObject
	res, err := s.s3Client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file from S3: %w", err)
	}
	defer res.Body.Close()

	// Read content from the result body
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read data from S3 object: %w", err)
	}

	// Return the file data and metadata
	return data, res.Metadata, nil
}

// --------------------------------------------------------------------

func (s *S3Service) GetSignedURL(path *string, fileName *string, expireSeconds int64) (string, error) {
	// Define the full key
	key := *path + *fileName

	// Create a pre-sign client
	presignClient := s3.NewPresignClient(s.s3Client)

	// Set up the GetObject input for pre-signing
	input := &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	}

	// Generate the pre-signed URL with expiration
	presignedReq, err := presignClient.PresignGetObject(context.TODO(), input, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expireSeconds) * time.Second
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %w", err)
	}

	if s.config.isLoggingEnabled {
		// util.LogTask("generated pre-signed URL for: "+key, "S3Service.GetSignedURL()", processID)
	}

	return presignedReq.URL, nil
}

// --------------------------------------------------------------------

func (s *S3Service) PostObject(path *string, fileName *string, data []byte, objectTitle string, processID string) error {
	// Define the full key
	key := *path + *fileName

	// Create a metadata to store real name of the file
	metadata := map[string]string{
		"title": objectTitle,
	}

	// Set up PutObject input with metadata
	input := &s3.PutObjectInput{
		Bucket:   &s.bucket,
		Key:      &key,
		Body:     bytes.NewReader(data),
		Metadata: metadata,
	}

	if s.config.isLoggingEnabled {
		// util.LogTask("uploading to S3: "+key, "S3Service.PostObject()", processID)
	}

	// Call S3 PutObject
	_, err := s.s3Client.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	if s.config.isLoggingEnabled {
		// util.LogTask("successfully uploaded to S3: "+key, "S3Service.PostObject()", processID)
	}

	return nil
}

// ----------------------------------------------------------------------------------------------------------------------------------------

func (s *S3Service) DeleteObject(path string, fileName string, processID string) error {
	// Define the full key for S3
	key := path + fileName

	// Set up the DeleteObject input
	input := &s3.DeleteObjectInput{
		Bucket: &s.bucket, // Replace with your bucket name variable
		Key:    &key,
	}

	// Call S3 DeleteObject
	_, err := s.s3Client.DeleteObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	// Logging to the console
	if s.config.isLoggingEnabled {
		// util.LogTask("successfully deleted from S3: "+key, "S3Service.DeleteObject()", processID)
	}

	return nil
}

// ----------------------------------------------------------------------------------------------------------------------------------------

func (s *S3Service) GetObjectHead(path string, fileName string, processID string) (*s3.HeadObjectOutput, error) {
	// Define the full key (path + fileName)
	key := path + fileName

	// Set up the HeadObject input
	input := &s3.HeadObjectInput{
		Bucket: &s.bucket, // S3 bucket name
		Key:    &key,      // S3 object key
	}

	// Call S3 HeadObject to retrieve metadata
	res, err := s.s3Client.HeadObject(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for file from S3: %w", err)
	}

	// Return the metadata result
	return res, nil
}
