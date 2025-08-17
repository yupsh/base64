package command

type WrapWidth int

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

type flags struct {
	Decode        DecodeFlag
	IgnoreGarbage IgnoreGarbageFlag
	Wrap          WrapFlag
	WrapWidth     WrapWidth
}

func (d DecodeFlag) Configure(flags *flags)        { flags.Decode = d }
func (i IgnoreGarbageFlag) Configure(flags *flags) { flags.IgnoreGarbage = i }
func (w WrapFlag) Configure(flags *flags)          { flags.Wrap = w }
func (w WrapWidth) Configure(flags *flags)         { flags.WrapWidth = w }
