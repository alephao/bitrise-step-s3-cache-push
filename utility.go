package main

import (
	"fmt"
	"os"
)

func CreateTempFolder(f func(tempFolderPath string)) {
	path, err := os.MkdirTemp("", "*")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	f(path)

	err = os.RemoveAll(path)

	if err != nil {
		fmt.Println(err)
	}
}

func GetEnvOrExit(key string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Println(fmt.Sprintf("Missing variable '%s'", key))
		os.Exit(1)
	}
	return value
}
