package command

import (
	"encoding/base64"
	"strings"

	gloo "github.com/gloo-foo/framework"
)

type command gloo.Inputs[gloo.File, flags]

func Base64(parameters ...any) gloo.Command {
	cmd := command(gloo.Initialize[gloo.File, flags](parameters...))
	if cmd.Flags.WrapWidth == 0 {
		cmd.Flags.WrapWidth = 76
	}
	return cmd
}

func (p command) Executor() gloo.CommandExecutor {
	return gloo.LineTransform(func(line string) (string, bool) {
		if bool(p.Flags.Decode) {
			// Decode mode
			input := line
			if !bool(p.Flags.IgnoreGarbage) {
				input = strings.TrimSpace(input)
			}

			decoded, decodeErr := base64.StdEncoding.DecodeString(input)
			if decodeErr != nil {
				if bool(p.Flags.IgnoreGarbage) {
					return "", false
				}
				return "", false
			}

			return string(decoded), true
		}

		// Encode mode
		encoded := base64.StdEncoding.EncodeToString([]byte(line))

		// Handle wrapping
		if bool(p.Flags.Wrap) {
			width := int(p.Flags.WrapWidth)
			if width == 0 {
				width = 76
			}

			var wrapped strings.Builder
			for i := 0; i < len(encoded); i += width {
				end := i + width
				if end > len(encoded) {
					end = len(encoded)
				}
				wrapped.WriteString(encoded[i:end])
				if end < len(encoded) {
					wrapped.WriteString("\n")
				}
			}
			return wrapped.String(), true
		}

		return encoded, true
	}).Executor()
}
