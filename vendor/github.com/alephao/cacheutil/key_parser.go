package cacheutil

import (
	"regexp"
	"strings"
)

type KeyParser struct {
	FunctionExecuter ICacheKeyFunctionExecuter
}

func NewKeyParser(executer ICacheKeyFunctionExecuter) KeyParser {
	return KeyParser{
		FunctionExecuter: executer,
	}
}

func parseFunctionName(str string) (f string, rest string) {
	trimmed := strings.TrimSpace(str)
	funcAndRest := strings.SplitN(trimmed, " ", 2)
	f = funcAndRest[0]

	if len(funcAndRest) > 1 {
		rest = funcAndRest[1]
	}

	return f, rest
}

func parseFunctionArgument(str string) (arg string, rest string) {
	trimmed := strings.TrimSpace(str)
	const (
		quote = 34
		space = 32
	)
	inQuote := false
	argChars := []rune{}
	for _, c := range trimmed {
		if c == quote {
			inQuote = !inQuote
		}

		if c == space && !inQuote {
			break
		}

		argChars = append(argChars, c)
	}

	argWithQuotes := string(argChars)
	rest = strings.TrimPrefix(trimmed, argWithQuotes)
	arg = strings.Trim(argWithQuotes, "\"")
	return arg, rest
}

func parseFunctionAndArgs(str string) []string {
	trimmed := strings.TrimSpace(str)

	f, rest := parseFunctionName(trimmed)

	result := []string{f}

	for len(rest) > 0 {
		arg, r := parseFunctionArgument(rest)
		rest = r
		result = append(result, arg)
	}

	return result
}

func ParseFunctions(key string) map[string][]string {
	curlyBracesRegex := regexp.MustCompile(`{{(.+?)}}`)

	allSubmatches := curlyBracesRegex.FindAllStringSubmatch(key, -1)

	dict := map[string][]string{}
	for _, match := range allSubmatches {
		fullMatch := match[0]
		funcAndArgs := parseFunctionAndArgs(match[1])
		dict[fullMatch] = funcAndArgs
	}
	return dict
}

func (kp *KeyParser) Parse(key string) (string, error) {
	fs := ParseFunctions(key)

	result := key
	for k, funcAndArgs := range fs {
		replacement, err := kp.FunctionExecuter.Execute(funcAndArgs)

		if err != nil {
			return "", err
		}

		result = strings.Replace(result, k, replacement, 1)
	}

	return result, nil
}
