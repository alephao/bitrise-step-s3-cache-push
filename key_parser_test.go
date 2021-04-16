package main

import "testing"

type ChecksumEngineMock struct {
	checksum string
}

func NewChecksumEngineMock(checksum string) *ChecksumEngineMock {
	return &ChecksumEngineMock{
		checksum: checksum,
	}
}

func (e ChecksumEngineMock) ChecksumForFile(file string) string {
	return e.checksum
}

func TestParseKey(t *testing.T) {
	engine := NewChecksumEngineMock("123abc")
	keyParser := NewKeyParser(engine)

	t.Run("simple key", func(t *testing.T) {
		key := "lock-file"
		parsed := keyParser.parse(key)
		expected := "lock-file"

		if parsed != expected {
			t.Errorf("Expected parsed key to be '%s' but got '%s'", expected, parsed)
		}
	})

	t.Run("key with single checksum sub", func(t *testing.T) {
		key := "lock-{{ checksum \"Lockfile\" }}"
		parsed := keyParser.parse(key)
		expected := "lock-123abc"

		if parsed != expected {
			t.Errorf("Expected parsed key to be '%s' but got '%s'", expected, parsed)
		}
	})

	t.Run("key with multiple checksum sub", func(t *testing.T) {
		key := "lock-{{ checksum \"Lockfile\" }}-blabla-{{ checksum \"Anotherfile\" }}"
		parsed := keyParser.parse(key)
		expected := "lock-123abc-blabla-123abc"

		if parsed != expected {
			t.Errorf("Expected parsed key to be '%s' but got '%s'", expected, parsed)
		}
	})
}
