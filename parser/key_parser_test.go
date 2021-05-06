package parser

import (
	"reflect"
	"testing"
)

func TestParseFunctions(t *testing.T) {
	t.Run("function with no arguments", func(t *testing.T) {
		key := "foo-{{ f }}"
		result := ParseFunctions(key)
		expected := map[string][]string{
			"{{ f }}": {"f"},
		}
		assertMapEqual(t, expected, result)
	})

	t.Run("function with no argument and no spacing", func(t *testing.T) {
		key := "foo-{{f}}"
		result := ParseFunctions(key)
		expected := map[string][]string{
			"{{f}}": {"f"},
		}
		assertMapEqual(t, expected, result)
	})

	t.Run("function with one argument", func(t *testing.T) {
		key := "foo-{{ f arg1 }}"
		result := ParseFunctions(key)
		expected := map[string][]string{
			"{{ f arg1 }}": {"f", "arg1"},
		}
		assertMapEqual(t, expected, result)
	})

	t.Run("function with one argument and multiple spaces in between", func(t *testing.T) {
		key := "foo-{{ f   arg1 }}"
		result := ParseFunctions(key)
		expected := map[string][]string{
			"{{ f   arg1 }}": {"f", "arg1"},
		}
		assertMapEqual(t, expected, result)
	})

	t.Run("function with one argument containing spaces in between quotes", func(t *testing.T) {
		key := "foo-{{ f \"my arg1\" }}"
		result := ParseFunctions(key)
		expected := map[string][]string{
			"{{ f \"my arg1\" }}": {"f", "my arg1"},
		}
		assertMapEqual(t, expected, result)
	})

	t.Run("function with multiple arguments", func(t *testing.T) {
		key := "foo-{{ f \"my arg1\" arg2 \"my arg 3\" }}"
		result := ParseFunctions(key)
		expected := map[string][]string{
			"{{ f \"my arg1\" arg2 \"my arg 3\" }}": {"f", "my arg1", "arg2", "my arg 3"},
		}
		assertMapEqual(t, expected, result)
	})
}

func assertMapEqual(t *testing.T, expected, result map[string][]string) {
	t.Helper()

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected:\n%+v\n\nBut got:\n%+v", expected, result)
	}
}

type CacheKeyFunctionExecuterMock struct {
	Result string
	Err    error
}

func NewCacheKeyFunctionExecuterMockSuccess(result string) CacheKeyFunctionExecuterMock {
	return CacheKeyFunctionExecuterMock{result, nil}
}

func NewCacheKeyFunctionExecuterMockErr(err error) CacheKeyFunctionExecuterMock {
	return CacheKeyFunctionExecuterMock{"", err}
}

func (e CacheKeyFunctionExecuterMock) Execute(funcAndArgs []string) (string, error) {
	return e.Result, e.Err
}

func TestParse(t *testing.T) {
	keyParser := NewKeyParser(
		NewCacheKeyFunctionExecuterMockSuccess("foo"),
	)

	result, err := keyParser.Parse("bar-{{ f }}-{{ g a b }}")
	expected := "bar-foo-foo"

	if err != nil {
		t.Errorf("expected error to be nil, but got '%s'", err.Error())
	}

	if result != expected {
		t.Errorf("expected %s, but got %s", expected, result)
	}
}
