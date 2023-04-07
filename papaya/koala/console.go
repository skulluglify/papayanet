package koala

import (
  "PapayaNet/papaya/koala/kio/leaf"
  "PapayaNet/papaya/koala/os/runtime"
  "fmt"
  "reflect"
  "strconv"
)

//goland:noinspection ALL
const (
  TypeInfo = iota
  TypeError
  TypeWarn
)

//goland:noinspection ALL
const (
  ColorBlack = iota
  ColorRed
  ColorGreen
  ColorYellow
  ColorBlue
  ColorPurple
  ColorCyan
  ColorWhite
)

//goland:noinspection ALL
const (
  StyleLight     = 0
  StyleBold      = 1
  StyleItalic    = 3
  StyleUnderline = 4
  StyleInvert    = 7
)

type KConsole struct {
  Platform  int
  Arch      int
  Colorful  bool
  Silent    bool
  Pad       int
  listeners []func(ntype int) error
  suspend   bool
}

// TODO: maybe added log file handling

type KConsoleImpl interface {
  Init()
  Main(event func() error)
  Listen(event func(ntype int) error)
  Dispatch(info int)
  Text(text string, color int, bgColor int, styles ...int) *leaf.KBuffer
  Print(info int, buffers ...leaf.KBufferImpl)
  Log(args ...any)
  Error(args ...any)
  Fatal(args ...any)
  Warn(args ...any)
  Space() string
  EOL() string
}

func KConsoleIsSupportColorful() bool {

  platform := runtime.KPlatformLook()

  switch platform {
  case runtime.PlatformLinux, runtime.PlatformDarwin:
    return true
  }

  return false
}

func KConsoleNew() KConsoleImpl {

  console := &KConsole{}
  console.Init()

  return console
}

func (c *KConsole) Init() {

  c.Platform = runtime.KPlatformLook()
  c.Arch = runtime.KArchLook()
  c.Colorful = KConsoleIsSupportColorful()
  c.Silent = false
  c.Pad = 1
}

func (c *KConsole) Main(event func() error) {

  if !c.Silent {

    if err := event(); err != nil {

      c.Error(err)
    }
  }
}

func (c *KConsole) Listen(event func(ntype int) error) {

  if !c.Silent {

    c.listeners = append(c.listeners, event)
  }
}

func (c *KConsole) Dispatch(info int) {

  // TODO: maybe prevent silent
  if !c.Silent {

    if !c.suspend {
      c.suspend = true
      for _, listener := range c.listeners {

        if err := listener(info); err != nil {

          c.Error(err)
        }
      }
      c.suspend = false
    }
  }
}

func (c *KConsole) Text(text string, color int, bgColor int, styles ...int) *leaf.KBuffer {

  if c.Colorful {

    var style string

    for _, v := range styles {

      style += strconv.Itoa(v) + ";"
    }

    return leaf.KMakeBuffer([]byte("\x1b" + "[" + style + strconv.Itoa(30+color) + ";" + strconv.Itoa(40+bgColor) + "m" + text + "\x1b[0m"))
  }

  return leaf.KMakeBuffer([]byte(text))
}

func (c *KConsole) Print(i int, buffers ...leaf.KBufferImpl) {

  t := KDateTimeNew()

  info := "UNKNOWN"
  infoColor := ColorCyan

  switch i {
  case TypeInfo:
    info = "INFO"
    infoColor = ColorCyan
  case TypeError:
    info = "ERROR"
    infoColor = ColorRed
  case TypeWarn:
    info = "WARN"
    infoColor = ColorYellow
  }

  tb := c.Text(
    "["+t.Ez()+"]",
    ColorGreen,
    ColorBlack,
    StyleLight)

  ib := c.Text(
    "["+info+"]",
    infoColor,
    ColorBlack,
    StyleBold)

  print(string(tb.ReadAll()), c.Space())
  print(string(ib.ReadAll()), c.Space())

  tb.Close()
  ib.Close()

  for _, buffer := range buffers {

    print(string(buffer.ReadAll()))
    print(c.Space())
    buffer.Close()
  }

  print(c.EOL())
}

func (c *KConsole) Log(args ...any) {

  if !c.Silent {

    c.Dispatch(TypeInfo)

    buffers := make([]leaf.KBufferImpl, 0)

    for _, v := range args {

      val := reflect.ValueOf(v)

      if val.IsValid() {

        ty := val.Type()

        switch ty.Kind() {

        //case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
        //  continue
        case reflect.Struct:

          if ty == leaf.KBufferType {

            buffers = append(buffers, v.(*leaf.KBuffer))
            continue
          }

        case reflect.Pointer:

          elem := ty.Elem()

          if elem == leaf.KBufferType {

            buffers = append(buffers, v.(*leaf.KBuffer))
            continue
          }
        }

        buffers = append(buffers, c.Text(fmt.Sprint(v),
          ColorWhite,
          ColorBlack,
          StyleLight,
        ))

        continue
      }

      buffers = append(buffers, c.Text("NULL",
        ColorCyan,
        ColorBlack,
        StyleLight,
      ))
    }

    c.Print(TypeInfo, buffers...)
  }
}

func (c *KConsole) Error(args ...any) {

  if !c.Silent {

    c.Dispatch(TypeError)

    buffers := make([]leaf.KBufferImpl, 0)

    for _, v := range args {

      val := reflect.ValueOf(v)

      if val.IsValid() {

        ty := val.Type()

        switch ty.Kind() {

        //case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
        //  continue
        case reflect.Struct:

          if ty == leaf.KBufferType {

            buffers = append(buffers, v.(*leaf.KBuffer))
            continue
          }

        case reflect.Pointer:

          elem := ty.Elem()

          if elem == leaf.KBufferType {

            buffers = append(buffers, v.(*leaf.KBuffer))
            continue
          }
        }

        buffers = append(buffers, c.Text(fmt.Sprint(v),
          ColorRed,
          ColorBlack,
          StyleLight,
        ))

        continue
      }

      buffers = append(buffers, c.Text("NULL",
        ColorCyan,
        ColorBlack,
        StyleLight,
      ))
    }

    c.Print(TypeError, buffers...)
  }
}

func (c *KConsole) Fatal(args ...any) {

  // link
  c.Error(args...)
}

func (c *KConsole) Warn(args ...any) {

  if !c.Silent {

    c.Dispatch(TypeWarn)

    buffers := make([]leaf.KBufferImpl, 0)

    for _, v := range args {

      val := reflect.ValueOf(v)

      if val.IsValid() {

        ty := val.Type()

        switch ty.Kind() {

        //case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
        //  continue
        case reflect.Struct:

          if ty == leaf.KBufferType {

            buffers = append(buffers, v.(*leaf.KBuffer))
            continue
          }

        case reflect.Pointer:

          elem := ty.Elem()

          if elem == leaf.KBufferType {

            buffers = append(buffers, v.(*leaf.KBuffer))
            continue
          }
        }

        buffers = append(buffers, c.Text(fmt.Sprint(v),
          ColorYellow,
          ColorBlack,
          StyleLight,
        ))

        continue
      }

      buffers = append(buffers, c.Text("NULL",
        ColorCyan,
        ColorBlack,
        StyleLight,
      ))
    }

    c.Print(TypeWarn, buffers...)
  }
}

func (c *KConsole) Space() string {

  pads := ""

  for i := 0; i < c.Pad; i++ {

    pads += " "
  }

  return pads
}

func (c *KConsole) EOL() string {

  return "\n"
}
