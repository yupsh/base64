package base64

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"
	localopt "github.com/yupsh/base64/opt"
)

// Flags represents the configuration options for the base64 command
type Flags = localopt.Flags

// Command implementation
type command opt.Inputs[string, Flags]

// Base64 creates a new base64 command with the given parameters
func Base64(parameters ...any) yup.Command {
	cmd := command(opt.Args[string, Flags](parameters...))
	// Set default wrap width
	if cmd.Flags.WrapWidth == 0 {
		cmd.Flags.WrapWidth = 76
	}
	return cmd
}

func (c command) Execute(ctx context.Context, input io.Reader, output, stderr io.Writer) error {
	if bool(c.Flags.Decode) {
		return c.decode(ctx, input, output, stderr)
	} else {
		return c.encode(ctx, input, output, stderr)
	}
}

func (c command) encode(ctx context.Context, input io.Reader, output, stderr io.Writer) error {
	return yup.ProcessFilesWithContext(
		ctx, c.Positional, input, output, stderr,
		yup.FileProcessorOptions{
			CommandName:     "base64",
			ContinueOnError: true,
		},
		func(ctx context.Context, source yup.InputSource, output io.Writer) error {
			return c.encodeSource(ctx, source.Reader, output)
		},
	)
}

func (c command) encodeSource(ctx context.Context, reader io.Reader, output io.Writer) error {
	// Check for cancellation before starting
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Read data in chunks to support cancellation for large files
	const chunkSize = 32 * 1024 // 32KB chunks
	var allData []byte
	buf := make([]byte, chunkSize)

	for {
		// Check for cancellation before each read
		if err := yup.CheckContextCancellation(ctx); err != nil {
			return err
		}

		n, err := reader.Read(buf)
		if n > 0 {
			allData = append(allData, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	encoded := base64.StdEncoding.EncodeToString(allData)

	if bool(c.Flags.Wrap) && int(c.Flags.WrapWidth) > 0 {
		wrapped, err := c.wrapStringWithContext(ctx, encoded, int(c.Flags.WrapWidth))
		if err != nil {
			return err
		}
		encoded = wrapped
	}

	fmt.Fprintln(output, encoded)
	return nil
}

func (c command) decode(ctx context.Context, input io.Reader, output, stderr io.Writer) error {
	return yup.ProcessFilesWithContext(
		ctx, c.Positional, input, output, stderr,
		yup.FileProcessorOptions{
			CommandName:     "base64",
			ContinueOnError: true,
		},
		func(ctx context.Context, source yup.InputSource, output io.Writer) error {
			return c.decodeSource(ctx, source.Reader, output, stderr)
		},
	)
}

func (c command) decodeSource(ctx context.Context, reader io.Reader, output io.Writer, stderr io.Writer) error {
	// Check for cancellation before starting
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)
	var encodedData strings.Builder

	for yup.ScanWithContext(ctx, scanner) {
		line := strings.TrimSpace(scanner.Text())
		if bool(c.Flags.IgnoreGarbage) {
			cleaned, err := c.removeNonBase64WithContext(ctx, line)
			if err != nil {
				return err
			}
			line = cleaned
		}
		encodedData.WriteString(line)
	}

	// Check if context was cancelled
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(stderr, "base64: %v\n", err)
		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedData.String())
	if err != nil {
		fmt.Fprintf(stderr, "base64: %v\n", err)
		return err
	}

	output.Write(decoded)
	return nil
}

func (c command) wrapString(s string, width int) string {
	if width <= 0 {
		return s
	}

	var result strings.Builder
	for i, char := range s {
		if i > 0 && i%width == 0 {
			result.WriteRune('\n')
		}
		result.WriteRune(char)
	}

	return result.String()
}

func (c command) wrapStringWithContext(ctx context.Context, s string, width int) (string, error) {
	if width <= 0 {
		return s, nil
	}

	var result strings.Builder
	for i, char := range s {
		// Check for cancellation every 1000 characters for efficiency
		if i%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return "", err
			}
		}

		if i > 0 && i%width == 0 {
			result.WriteRune('\n')
		}
		result.WriteRune(char)
	}

	return result.String(), nil
}

func (c command) removeNonBase64(s string) string {
	var result strings.Builder
	for _, char := range s {
		if (char >= 'A' && char <= 'Z') ||
			(char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			char == '+' || char == '/' || char == '=' {
			result.WriteRune(char)
		}
	}
	return result.String()
}

func (c command) removeNonBase64WithContext(ctx context.Context, s string) (string, error) {
	var result strings.Builder
	for i, char := range s {
		// Check for cancellation every 1000 characters for efficiency
		if i%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return "", err
			}
		}

		if (char >= 'A' && char <= 'Z') ||
			(char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			char == '+' || char == '/' || char == '=' {
			result.WriteRune(char)
		}
	}
	return result.String(), nil
}

func (c command) String() string {
	return fmt.Sprintf("base64 %v", c.Positional)
}
