package base64_test

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"

	yup "github.com/yupsh/framework"

	"github.com/yupsh/base64"
	"github.com/yupsh/base64/opt"
)

// Example tests (basic functionality)
func ExampleBase64() {
	ctx := context.Background()
	input := strings.NewReader("Hello World")

	cmd := base64.Base64()
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output: SGVsbG8gV29ybGQ=
}

func ExampleBase64_decode() {
	ctx := context.Background()
	input := strings.NewReader("SGVsbG8gV29ybGQ=")

	cmd := base64.Base64(opt.Decode)
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output: Hello World
}

func ExampleBase64_withWrapping() {
	ctx := context.Background()
	input := strings.NewReader("This is a longer string that will be wrapped")

	cmd := base64.Base64(opt.WrapWidth(20))
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output: VGhpcyBpcyBhIGxvbmdl
	// ciBzdHJpbmcgdGhhdCB3
	// aWxsIGJlIHdyYXBwZWQ=
}

// Comprehensive functionality tests
func TestBase64_BasicEncode(t *testing.T) {
	ctx := context.Background()
	input := strings.NewReader("test data")
	var output, stderr bytes.Buffer

	cmd := base64.Base64()
	err := cmd.Execute(ctx, input, &output, &stderr)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "dGVzdCBkYXRh"
	if got := strings.TrimSpace(output.String()); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestBase64_BasicDecode(t *testing.T) {
	ctx := context.Background()
	input := strings.NewReader("dGVzdCBkYXRh")
	var output, stderr bytes.Buffer

	cmd := base64.Base64(opt.Decode)
	err := cmd.Execute(ctx, input, &output, &stderr)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "test data"
	if got := strings.TrimSpace(output.String()); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

// Context cancellation tests (CRITICAL for large files)
func TestBase64_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create large input that would take time to process
	largeInput := strings.Repeat("This is test data that will be repeated many times. ", 10000)
	input := strings.NewReader(largeInput)

	var output, stderr bytes.Buffer
	cmd := base64.Base64()

	// Start processing in goroutine
	done := make(chan error, 1)
	go func() {
		done <- cmd.Execute(ctx, input, &output, &stderr)
	}()

	// Cancel after short time
	time.Sleep(10 * time.Millisecond)
	cancel()

	// Verify cancellation behavior
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Error("Command did not respond to cancellation within timeout")
	}
}

func TestBase64_DecodeContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create large base64 input
	largeData := strings.Repeat("This is test data. ", 10000)
	_ = largeData // Use the variable to avoid unused error
	input := strings.NewReader(strings.Repeat("dGVzdA==", 10000))

	var output, stderr bytes.Buffer
	cmd := base64.Base64(opt.Decode)

	done := make(chan error, 1)
	go func() {
		done <- cmd.Execute(ctx, input, &output, &stderr)
	}()

	// Cancel quickly
	time.Sleep(5 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Error("Decode did not respond to cancellation")
	}
}

// Error condition tests
func TestBase64_ErrorConditions(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		decode    bool
		wantError bool
	}{
		{"empty input encode", "", false, false},
		{"empty input decode", "", true, false},
		{"invalid base64", "invalid base64!", true, true},
		{"partial base64", "SGVsbG8", true, true},
		{"valid base64", "SGVsbG8=", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			input := strings.NewReader(tt.input)
			var output, stderr bytes.Buffer

			var cmd yup.Command
			if tt.decode {
				cmd = base64.Base64(opt.Decode)
			} else {
				cmd = base64.Base64()
			}

			err := cmd.Execute(ctx, input, &output, &stderr)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Flag combination tests
func TestBase64_WrapWidth(t *testing.T) {
	tests := []struct {
		name      string
		wrapWidth int
		input     string
		wantLines int
	}{
		{"no wrap", 0, "Hello World Test Data", 1},
		{"wrap 10", 10, "Hello World Test Data", 3},
		{"wrap 20", 20, "Hello World Test Data", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			input := strings.NewReader(tt.input)
			var output, stderr bytes.Buffer

			var cmd yup.Command
			if tt.wrapWidth > 0 {
				cmd = base64.Base64(opt.WrapWidth(tt.wrapWidth))
			} else {
				cmd = base64.Base64()
			}

			err := cmd.Execute(ctx, input, &output, &stderr)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			lines := strings.Split(strings.TrimSpace(output.String()), "\n")
			if len(lines) != tt.wantLines {
				t.Errorf("Expected %d lines, got %d: %v", tt.wantLines, len(lines), lines)
			}
		})
	}
}

// Performance benchmarks
func BenchmarkBase64_Encode(b *testing.B) {
	ctx := context.Background()
	testData := strings.Repeat("test data ", 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := strings.NewReader(testData)
		var output, stderr bytes.Buffer
		cmd := base64.Base64()

		err := cmd.Execute(ctx, input, &output, &stderr)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkBase64_Decode(b *testing.B) {
	ctx := context.Background()
	// Pre-encoded test data
	testData := "dGVzdCBkYXRhIA==" // "test data " in base64, repeated would be needed

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := strings.NewReader(testData)
		var output, stderr bytes.Buffer
		cmd := base64.Base64(opt.Decode)

		err := cmd.Execute(ctx, input, &output, &stderr)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}
