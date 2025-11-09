package command_test

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/base64"
)

// Test basic encoding
func TestBase64_Encode(t *testing.T) {
	result := run.Command(command.Base64()).
		WithStdinLines("hello").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"aGVsbG8="})
}

// Test basic decoding
func TestBase64_Decode(t *testing.T) {
	result := run.Command(command.Base64(command.Decode)).
		WithStdinLines("aGVsbG8=").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"hello"})
}

// Test multiple lines
func TestBase64_MultipleLines(t *testing.T) {
	result := run.Command(command.Base64()).
		WithStdinLines("line1", "line2").Run()
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 2)
}

// Test empty input
func TestBase64_EmptyInput(t *testing.T) {
	result := run.Quick(command.Base64())
	assertion.NoError(t, result.Err)
	assertion.Empty(t, result.Stdout)
}

// Test empty line
func TestBase64_EmptyLine(t *testing.T) {
	result := run.Command(command.Base64()).
		WithStdinLines("").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{""})
}

// Test wrap flag
func TestBase64_Wrap(t *testing.T) {
	result := run.Command(command.Base64(command.Wrap)).
		WithStdinLines("a very long string that will be wrapped").Run()
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 1)
}

// Test custom wrap width
func TestBase64_WrapWidth(t *testing.T) {
	result := run.Command(command.Base64(command.Wrap, command.WrapWidth(10))).
		WithStdinLines("hello world").Run()
	assertion.NoError(t, result.Err)
	// Should wrap at 10 characters
	assertion.Count(t, result.Stdout, 2)
}

// Test ignore garbage
func TestBase64_IgnoreGarbage(t *testing.T) {
	result := run.Command(command.Base64(command.Decode, command.IgnoreGarbage)).
		WithStdinLines("invalid!!!").Run()
	assertion.NoError(t, result.Err)
	// Invalid input with ignore garbage should not output
	assertion.Empty(t, result.Stdout)
}

// Test decode with whitespace
func TestBase64_DecodeWithWhitespace(t *testing.T) {
	result := run.Command(command.Base64(command.Decode)).
		WithStdinLines("  aGVsbG8=  ").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"hello"})
}

// Test output error
func TestBase64_OutputError(t *testing.T) {
	result := run.Command(command.Base64()).
		WithStdinLines("test").
		WithStdoutError(errors.New("write failed")).Run()
	assertion.ErrorContains(t, result.Err, "write failed")
}

// Test input error
func TestBase64_InputError(t *testing.T) {
	result := run.Command(command.Base64()).
		WithStdinError(errors.New("read failed")).Run()
	assertion.ErrorContains(t, result.Err, "read failed")
}

// Test various encodings
func TestBase64_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"hello", "hello", "aGVsbG8="},
		{"world", "world", "d29ybGQ="},
		{"123", "123", "MTIz"},
		{"abc", "abc", "YWJj"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := run.Command(command.Base64()).
				WithStdinLines(tt.input).Run()
			assertion.NoError(t, result.Err)
			assertion.Lines(t, result.Stdout, []string{tt.expected})
		})
	}
}

// Test decode invalid
func TestBase64_DecodeInvalid(t *testing.T) {
	result := run.Command(command.Base64(command.Decode)).
		WithStdinLines("not-valid-base64!!!").Run()
	assertion.NoError(t, result.Err)
	// Invalid decode should not output
	assertion.Empty(t, result.Stdout)
}

// Test unicode
func TestBase64_Unicode(t *testing.T) {
	result := run.Command(command.Base64()).
		WithStdinLines("日本語").Run()
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 1)

	// Decode it back
	encoded := result.Stdout[0]
	result2 := run.Command(command.Base64(command.Decode)).
		WithStdinLines(encoded).Run()
	assertion.NoError(t, result2.Err)
	assertion.Lines(t, result2.Stdout, []string{"日本語"})
}

