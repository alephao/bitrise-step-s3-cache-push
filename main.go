package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alephao/bitrise-step-s3-cache-push/parser"
	"github.com/mholt/archiver"
)

const (
	BITRISE_GIT_BRANCH       = "BITRISE_GIT_BRANCH"
	BITRISE_OSX_STACK_REV_ID = "BITRISE_OSX_STACK_REV_ID"

	CACHE_AWS_ACCESS_KEY_ID     = "cache_aws_access_key_id"
	CACHE_AWS_SECRET_ACCESS_KEY = "cache_aws_secret_access_key"
	CACHE_AWS_REGION            = "cache_aws_region"
	CACHE_BUCKET_NAME           = "cache_bucket_name"
	CACHE_KEY                   = "cache_key"
	CACHE_PATH                  = "cache_path"
)

func generateBucketKey(cacheKey string) (string, error) {
	branch := os.Getenv(BITRISE_GIT_BRANCH)
	stackrev := os.Getenv(BITRISE_OSX_STACK_REV_ID)
	functionExecuter := parser.NewCacheKeyFunctionExecuter(branch, stackrev)
	keyParser := parser.NewKeyParser(&functionExecuter)
	return keyParser.Parse(cacheKey)
}

func main() {
	awsAccessKeyId := GetEnvOrExit(CACHE_AWS_ACCESS_KEY_ID)
	awsSecretAccessKey := GetEnvOrExit(CACHE_AWS_SECRET_ACCESS_KEY)
	awsRegion := GetEnvOrExit(CACHE_AWS_REGION)
	bucketName := GetEnvOrExit(CACHE_BUCKET_NAME)
	cacheKey := GetEnvOrExit(CACHE_KEY)
	cachePath := GetEnvOrExit(CACHE_PATH)

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
