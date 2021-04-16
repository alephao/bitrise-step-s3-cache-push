package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

type ChecksumEngine interface {
	ChecksumForFile(path string) string
}

type FileChecksumEngine struct{}

func NewFileChecksumEngine() FileChecksumEngine {
	return FileChecksumEngine{}
}

func (e *FileChecksumEngine) ChecksumForFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
