package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alephao/bitrise-step-s3-cache-push/parser"
	"github.com/mholt/archiver"
)

const (
	BITRISE_GIT_BRANCH = "BITRISE_GIT_BRANCH"
)

func generateBucketKey(cacheKey string) (string, error) {
	branch := os.Getenv(BITRISE_GIT_BRANCH)
	functionExecuter := parser.NewCacheKeyFunctionExecuter(branch)
	keyParser := parser.NewKeyParser(&functionExecuter)
	return keyParser.Parse(cacheKey)
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
		bucketKey, err := generateBucketKey(cacheKey)

		if err != nil {
			log.Printf("Failed to parse cache key '%s'\n", cacheKey)
			log.Printf("Error: %s\n", err.Error())
			failed = true
			return
		}

		log.Printf("Checking if cache exists for key '%s'\n", bucketKey)
		cacheExists := s3.CacheExists(bucketKey)

		if cacheExists {
			log.Println("Cache found! Skiping...")
			return
		}

		log.Println("Cache not found, trying to compress the folder.")

		outputPath := fmt.Sprintf("%s/%s.zip", tempFolderPath, bucketKey)
		err = archiver.Archive([]string{cachePath}, outputPath)

		if err != nil {
			log.Printf("Failed to compress '%s'\n", cachePath)
			log.Printf("Error: %s\n", err.Error())
			failed = true
			return
		}

		log.Println("Compression was successful, trying to upload to aws.")

		err = s3.UploadToAws(
			bucketKey,
			outputPath,
		)

		if err != nil {
			log.Printf("Failed to upload! Failing gracefully. Error: %s\n", err)
			return
		}

		log.Println("Upload was successful!")
	})

	if failed {
		os.Exit(1)
	}

	os.Exit(0)
}
