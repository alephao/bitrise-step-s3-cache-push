package main

import (
	"regexp"
	"strings"
)

type FileMatch struct {
	FullMatch string
	FileName  string
}

func NewPair(fullMatch, fileName string) FileMatch {
	return FileMatch{
		FullMatch: fullMatch,
		FileName:  fileName,
	}
}

type KeyParser struct {
	Checksum ChecksumEngine
}

func NewKeyParser(checksum ChecksumEngine) *KeyParser {
	return &KeyParser{
		Checksum: checksum,
	}
}

func (p *KeyParser) parse(key string) string {
	pairs := p.parseChecksums(key)

	newKey := key

	for _, pair := range pairs {
		cksum := p.Checksum.ChecksumForFile(pair.FileName)
		newKey = strings.Replace(newKey, pair.FullMatch, cksum, 1)
	}

	return newKey
}

func (p *KeyParser) parseChecksums(key string) []FileMatch {
	re := regexp.MustCompile(`{{\s?checksum\s\"(.+?)\"\s?}}`)

	var allPairs []FileMatch
	submatchall := re.FindAllStringSubmatch(key, -1)
	for _, match := range submatchall {
		pair := NewPair(
			match[0],
			match[1],
		)
		allPairs = append(allPairs, pair)
	}

	return allPairs
}
