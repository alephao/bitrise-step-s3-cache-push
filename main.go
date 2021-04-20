package main

import (
	"fmt"
	"os"

	"github.com/mholt/archiver"
)

func GenerateBucketKey(cacheKey string) string {
	checksumEngine := NewFileChecksumEngine()
	keyParser := NewKeyParser(&checksumEngine)
	return keyParser.parse(cacheKey)
}

func main() {
	awsAccessKeyId := GetEnvOrExit("aws_access_key_id")
	awsSecretAccessKey := GetEnvOrExit("aws_secret_access_key")
	awsRegion := GetEnvOrExit("aws_region")
	bucketName := GetEnvOrExit("bucket_name")
	cacheKey := GetEnvOrExit("key")
	cachePath := GetEnvOrExit("path")

	failed := false

	CreateTempFolder(func(tempFolderPath string) {
		s3 := NewAwsS3(
			awsRegion,
			awsAccessKeyId,
			awsSecretAccessKey,
			bucketName,
		)
		bucketKey := GenerateBucketKey(cacheKey)

		fmt.Printf("Checking if cache exists for key '%s'\n", bucketKey)
		cacheExists := s3.CacheExists(bucketKey)

		if cacheExists {
			fmt.Println("Cache found! Skiping...")
			return
		}

		fmt.Println("Cache not found, trying to compress the folder upload.")

		outputPath := fmt.Sprintf("%s/%s.tar.gz", tempFolderPath, bucketKey)
		err := archiver.Archive([]string{cachePath}, outputPath)

		if err != nil {
			fmt.Printf("Failed to compress '%s'\n", cachePath)
			failed = true
			return
		}

		fmt.Println("Compression was successful, trying to upload to aws.")

		err = s3.UploadToAws(
			bucketKey,
			outputPath,
		)

		if err != nil {
			fmt.Printf("Failed to upload! Failing gracefully. Error: %s\n", err)
			return
		}

		fmt.Println("Upload was successful!")
	})

	if failed {
		os.Exit(1)
	}

	os.Exit(0)
}