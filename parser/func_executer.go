package parser

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

type ICacheKeyFunctionExecuter interface {
	Execute([]string) (string, error)
}

type CacheKeyFunctionExecuter struct {
	CurrentBranch        string
	BitriseOsxStackRevId string
}

func NewCacheKeyFunctionExecuter(branch, bitriseOsxStackRevId string) CacheKeyFunctionExecuter {
	return CacheKeyFunctionExecuter{
		CurrentBranch:        branch,
		BitriseOsxStackRevId: bitriseOsxStackRevId,
	}
}

func (e *CacheKeyFunctionExecuter) branch(args []string) (string, error) {
	if len(args) > 0 {
		return "", fmt.Errorf("the branch function doesn't accept any args")
	}

	if e.CurrentBranch == "" {
		return "", fmt.Errorf("no branch available")
	}

	return e.CurrentBranch, nil
}

func checksumForFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (e *CacheKeyFunctionExecuter) checksum(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("the checksum function only accepts 1 argument, but got %d: %s", len(args), strings.Join(args, ", "))
	}

	filePath := args[0]

	return checksumForFile(filePath)
}

func (e *CacheKeyFunctionExecuter) stackrev(args []string) (string, error) {
	if len(args) > 0 {
		return "", fmt.Errorf("the branch function doesn't accept any args")
	}

	if e.BitriseOsxStackRevId == "" {
		return "", fmt.Errorf("env var BITRISE_OSX_STACK_REV_ID is not available")
	}

	return e.BitriseOsxStackRevId, nil
}

func (e *CacheKeyFunctionExecuter) Execute(funcAndArgs []string) (string, error) {
	f := funcAndArgs[0]
	args := funcAndArgs[1:]

	switch f {
	case "branch":
		return e.branch(args)
	case "checksum":
		return e.checksum(args)
	case "stackrev":
		return e.stackrev(args)
	}

	availableFunctions := []string{"branch", "checksum", "stackrev"}

	return "", fmt.Errorf("unknown function named '%s'. The available functions are: %s", f, strings.Join(availableFunctions, ", "))
}
