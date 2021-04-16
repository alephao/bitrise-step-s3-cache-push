package main

import (
	"os"
	"testing"
)

var testFileName = "Lockfile"

func TestFileChecksumEngine(t *testing.T) {
	CreateTempFolder(func(tempFolderPath string) {
		filePath := tempFolderPath + "/" + testFileName
		checksumEngine := NewFileChecksumEngine()

		CreateFile(filePath, "Hello World")

		checksum := checksumEngine.ChecksumForFile(filePath)
		expected := "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"

		if checksum != expected {
			t.Errorf("Expected checksum to be '%s' but got '%s'", expected, checksum)
		}
	})
}

func CreateFile(path, content string) {
	f, _ := os.Create(path)

	_, err := f.WriteString(content)
	if err != nil {
		f.Close()
		return
	}

	err = f.Close()
	if err != nil {
		return
	}
}
