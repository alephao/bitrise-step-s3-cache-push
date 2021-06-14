package parser

import "testing"

func TestCacheKeyFunctionExecuter(t *testing.T) {
	executer := NewCacheKeyFunctionExecuter("myBranch", "stackrevid")

	t.Run("branch function", func(t *testing.T) {
		result, err := executer.Execute([]string{"branch"})
		expected := "myBranch"

		if err != nil {
			t.Errorf("expected error to be nil, but got '%s'", err.Error())
		}

		if result != expected {
			t.Errorf("expected '%s' but got '%s'", expected, result)
		}
	})

	t.Run("branch function with argument", func(t *testing.T) {
		result, err := executer.Execute([]string{"branch", "arg 1"})
		expectedErr := "the branch function doesn't accept any args"

		if result != "" {
			t.Errorf("result to be empty, but got '%s'", result)
		}

		if err.Error() != expectedErr {
			t.Errorf("expected '%s' but got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("stackrev function", func(t *testing.T) {
		result, err := executer.Execute([]string{"stackrev"})
		expected := "stackrevid"

		if err != nil {
			t.Errorf("expected error to be nil, but got '%s'", err.Error())
		}

		if result != expected {
			t.Errorf("expected '%s' but got '%s'", expected, result)
		}
	})

}
