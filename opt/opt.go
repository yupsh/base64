package opt

// Custom types for parameters
type WrapWidth int

// Boolean flag types with constants
type DecodeFlag bool

const (
	Decode   DecodeFlag = true
	NoDecode DecodeFlag = false
)

type IgnoreGarbageFlag bool

const (
	IgnoreGarbage   IgnoreGarbageFlag = true
	NoIgnoreGarbage IgnoreGarbageFlag = false
)

type WrapFlag bool

const (
	Wrap   WrapFlag = true
	NoWrap WrapFlag = false
)

// Flags represents the configuration options for the base64 command
type Flags struct {
	Decode        DecodeFlag        // Decode data (-d)
	IgnoreGarbage IgnoreGarbageFlag // Ignore non-alphabet characters when decoding (-i)
	Wrap          WrapFlag          // Wrap encoded lines after COLS character (default 76)
	WrapWidth     WrapWidth         // Wrap width (-w)
}

// Configure methods for the opt system
func (d DecodeFlag) Configure(flags *Flags)        { flags.Decode = d }
func (i IgnoreGarbageFlag) Configure(flags *Flags) { flags.IgnoreGarbage = i }
func (w WrapFlag) Configure(flags *Flags)          { flags.Wrap = w }
func (w WrapWidth) Configure(flags *Flags)         { flags.WrapWidth = w }
