package utils

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

const (
	PN_CONSOLE_TYPE_INFO = iota
	PN_CONSOLE_TYPE_ERROR
	PN_CONSOLE_TYPE_WARN
)

const (
	PN_CONSOLE_COLOR_BLACK = iota
	PN_CONSOLE_COLOR_RED
	PN_CONSOLE_COLOR_GREEN
	PN_CONSOLE_COLOR_YELLOW
	PN_CONSOLE_COLOR_BLUE
	PN_CONSOLE_COLOR_PURPLE
	PN_CONSOLE_COLOR_CYAN
	PN_CONSOLE_COLOR_WHITE
)

const (
	PN_CONSOLE_STYLE_LIGHT     = 0
	PN_CONSOLE_STYLE_BOLD      = 1
	PN_CONSOLE_STYLE_ITALIC    = 3
	PN_CONSOLE_STYLE_UNDERLINE = 4
	PN_CONSOLE_STYLE_INVERT    = 7
)

const (
	PN_PLATFORM_UNKNOWN = iota
	PN_PLATFORM_ANDROID
	PN_PLATFORM_DARWIN
	PN_PLATFORM_DRAGONFLY
	PN_PLATFORM_FREEBSD
	PN_PLATFORM_LINUX
	PN_PLATFORM_NACL
	PN_PLATFORM_NETBSD
	PN_PLATFORM_OPENBSD
	PN_PLATFORM_PLAN9
	PN_PLATFORM_SOLARIS
	PN_PLATFORM_WINDOWS
)

const (
	PN_ARCH_UNKNOWN = iota
	PN_ARCH_386
	PN_ARCH_AMD64
	PN_ARCH_AMD64P32
	PN_ARCH_ARM
	PN_ARCH_ARMBE
	PN_ARCH_ARM64
	PN_ARCH_ARM64BE
	PN_ARCH_PPC64
	PN_ARCH_PPC64LE
	PN_ARCH_MIPS
	PN_ARCH_MIPSLE
	PN_ARCH_MIPS64
	PN_ARCH_MIPS64LE
	PN_ARCH_MIPS64P32
	PN_ARCH_MIPS64P32LE
	PN_ARCH_PPC
	PN_ARCH_RISCV64
	PN_ARCH_S390X
)

func PnPlatformLook() int {

	switch runtime.GOOS {
	case "android":
		return PN_PLATFORM_ANDROID
	case "darwin":
		return PN_PLATFORM_DARWIN
	case "dragonfly":
		return PN_PLATFORM_DRAGONFLY
	case "freebsd":
		return PN_PLATFORM_FREEBSD
	case "linux":
		return PN_PLATFORM_LINUX
	case "nacl":
		return PN_PLATFORM_NACL
	case "netbsd":
		return PN_PLATFORM_NETBSD
	case "openbsd":
		return PN_PLATFORM_OPENBSD
	case "plan9":
		return PN_PLATFORM_PLAN9
	case "solaris":
		return PN_PLATFORM_SOLARIS
	case "windows":
		return PN_PLATFORM_WINDOWS
	}

	return PN_PLATFORM_UNKNOWN
}

func PnArchLook() int {

	switch runtime.GOARCH {
	case "386":
		return PN_ARCH_386
	case "amd64":
		return PN_ARCH_AMD64
	case "amd64p32":
		return PN_ARCH_AMD64P32
	case "arm":
		return PN_ARCH_ARM
	case "armbe":
		return PN_ARCH_ARMBE
	case "arm64":
		return PN_ARCH_ARM64
	case "arm64be":
		return PN_ARCH_ARM64BE
	case "ppc64":
		return PN_ARCH_PPC64
	case "ppc64le":
		return PN_ARCH_PPC64LE
	case "mips":
		return PN_ARCH_MIPS
	case "mipsle":
		return PN_ARCH_MIPSLE
	case "mips64":
		return PN_ARCH_MIPS64
	case "mips64le":
		return PN_ARCH_MIPS64LE
	case "mips64p32":
		return PN_ARCH_MIPS64P32
	case "mips64p32le":
		return PN_ARCH_MIPS64P32LE
	case "ppc":
		return PN_ARCH_PPC
	case "riscv64":
		return PN_ARCH_RISCV64
	case "s390x":
		return PN_ARCH_S390X
	}

	return PN_ARCH_UNKNOWN
}

type PnConsole struct {
	Platform  int
	Arch      int
	Colorful  bool
	Silent    bool
	Pad       int
	listeners []func(info int) error
	suspend   bool
}

// TODO: maybe added log file handling

type PnConsoleImpl interface {
	Init() error
	Main(event func() error)
	Listen(event func(info int) error)
	Text(text string, color int, style int) *PnBuffer
	Print(info int, buffers ...PnBuffer)
	Log(args ...any)
	Error(args ...any)
	Warn(args ...any)
	Space() string
	EOL() string
	Dispatch(info int)
}

func PnConsoleIsSupportColorful() bool {

	platform := PnPlatformLook()

	switch platform {
	case PN_PLATFORM_LINUX, PN_PLATFORM_DARWIN:
		return true
	}

	return false
}

func PnMakeConsole() *PnConsole {

	console := &PnConsole{}
	err := console.Init()
	if err != nil {
		return nil
	}
	return console
}

func (console *PnConsole) Init() error {

	console.Platform = PnPlatformLook()
	console.Arch = PnArchLook()
	console.Colorful = PnConsoleIsSupportColorful()
	console.Silent = false
	console.Pad = 1

	return nil
}

func (console *PnConsole) Main(event func() error) {

	if !console.Silent {

		if err := event(); err != nil {

			console.Error(err)
		}
	}
}

func (console *PnConsole) Listen(event func(info int) error) {

	if !console.Silent {

		console.listeners = append(console.listeners, event)
	}
}

func (console *PnConsole) Dispatch(info int) {

	// TODO: maybe prevent silent
	if !console.Silent {

		if !console.suspend {
			console.suspend = true
			for _, listener := range console.listeners {

				if err := listener(info); err != nil {

					console.Error(err)
				}
			}
			console.suspend = false
		}
	}
}

func (console *PnConsole) Text(text string, color int, bgColor int, styles ...int) *PnBuffer {

	if console.Colorful {

		var style string

		for _, v := range styles {

			style += strconv.Itoa(v) + ";"
		}

		return PnMakeBuffer([]byte("\x1b" + "[" + style + strconv.Itoa(30+color) + ";" + strconv.Itoa(40+bgColor) + "m" + text + "\x1b[0m"))
	}

	return PnMakeBuffer([]byte(text))
}

func (console *PnConsole) Print(i int, buffers ...PnBufferImpl) {

	t := PnTimeNow()

	info := "UNKNOWN"
	infoColor := PN_CONSOLE_COLOR_CYAN

	switch i {
	case PN_CONSOLE_TYPE_INFO:
		info = "INFO"
		infoColor = PN_CONSOLE_COLOR_CYAN
	case PN_CONSOLE_TYPE_ERROR:
		info = "ERROR"
		infoColor = PN_CONSOLE_COLOR_RED
	case PN_CONSOLE_TYPE_WARN:
		info = "WARN"
		infoColor = PN_CONSOLE_COLOR_YELLOW
	}

	tb := console.Text(
		"["+t.Simple()+"]",
		PN_CONSOLE_COLOR_GREEN,
		PN_CONSOLE_COLOR_BLACK,
		PN_CONSOLE_STYLE_LIGHT)

	ib := console.Text(
		"["+info+"]",
		infoColor,
		PN_CONSOLE_COLOR_BLACK,
		PN_CONSOLE_STYLE_BOLD)

	print(string(tb.ReadAll()), console.Space())
	print(string(ib.ReadAll()), console.Space())

	tb.Close()
	ib.Close()

	for _, buffer := range buffers {

		print(string(buffer.ReadAll()))
		print(console.Space())
		buffer.Close()
	}

	print(console.EOL())
}

func (console *PnConsole) Log(args ...any) {

	if !console.Silent {

		console.Dispatch(PN_CONSOLE_TYPE_INFO)

		buffers := make([]PnBufferImpl, 0)

		for _, v := range args {

			typeOf := reflect.TypeOf(v)

			switch typeOf.Kind() {

			//case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
			//	continue
			case reflect.Struct:

				if typeOf == PnBufferType {

					buffers = append(buffers, v.(*PnBuffer))
					continue
				}

			case reflect.Pointer:

				elem := typeOf.Elem()

				if elem == PnBufferType {

					buffers = append(buffers, v.(*PnBuffer))
					continue
				}
			}

			buffers = append(buffers, console.Text(fmt.Sprint(v),
				PN_CONSOLE_COLOR_WHITE,
				PN_CONSOLE_COLOR_BLACK,
				PN_CONSOLE_STYLE_LIGHT,
			))
		}

		console.Print(PN_CONSOLE_TYPE_INFO, buffers...)
	}
}

func (console *PnConsole) Error(args ...any) {

	if !console.Silent {

		console.Dispatch(PN_CONSOLE_TYPE_ERROR)

		buffers := make([]PnBufferImpl, 0)

		for _, v := range args {

			typeOf := reflect.TypeOf(v)

			switch typeOf.Kind() {

			//case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
			//	continue
			case reflect.Struct:

				if typeOf == PnBufferType {

					buffers = append(buffers, v.(*PnBuffer))
					continue
				}

			case reflect.Pointer:

				elem := typeOf.Elem()

				if elem == PnBufferType {

					buffers = append(buffers, v.(*PnBuffer))
					continue
				}
			}

			buffers = append(buffers, console.Text(fmt.Sprint(v),
				PN_CONSOLE_COLOR_RED,
				PN_CONSOLE_COLOR_BLACK,
				PN_CONSOLE_STYLE_BOLD,
			))
		}

		console.Print(PN_CONSOLE_TYPE_ERROR, buffers...)
	}
}

func (console *PnConsole) Warn(args ...any) {

	if !console.Silent {

		console.Dispatch(PN_CONSOLE_TYPE_WARN)

		buffers := make([]PnBufferImpl, 0)

		for _, v := range args {

			typeOf := reflect.TypeOf(v)

			switch typeOf.Kind() {

			//case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
			//	continue
			case reflect.Struct:

				if typeOf == PnBufferType {

					buffers = append(buffers, v.(*PnBuffer))
					continue
				}

			case reflect.Pointer:

				elem := typeOf.Elem()

				if elem == PnBufferType {

					buffers = append(buffers, v.(*PnBuffer))
					continue
				}
			}

			buffers = append(buffers, console.Text(fmt.Sprint(v),
				PN_CONSOLE_COLOR_YELLOW,
				PN_CONSOLE_COLOR_BLACK,
				PN_CONSOLE_STYLE_LIGHT,
			))
		}

		console.Print(PN_CONSOLE_TYPE_WARN, buffers...)
	}
}

func (console *PnConsole) Space() string {

	pads := ""

	for i := 0; i < console.Pad; i++ {

		pads += " "
	}

	return pads
}

func (console *PnConsole) EOL() string {

	return "\n"
}
