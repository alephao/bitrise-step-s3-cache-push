package main

import (
	"fmt"
	"log"
	"os"
)

func CreateTempFolder(f func(tempFolderPath string)) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalln(err.Error())
	}

	path := fmt.Sprintf("%s/bitrise-s3-step-push-tmp", homeDir)

	err = os.MkdirAll(path, os.ModePerm)

	if err != nil {
		log.Fatalln(err.Error())
	}

	f(path)

	err = os.RemoveAll(path)

	if err != nil {
		log.Printf("failed to remove temp folder '%s', you'll have remove it manually\nError: %s\n", path, err.Error())
	}
}

func GetEnvOrExit(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing environment variable '%s", key)
	}
	return value
}
