package runtime

import "runtime"

//goland:noinspection ALL
const (
  PlatformUnknown = iota
  PlatformAndroid
  PlatformDarwin
  PlatformDragonfly
  PlatformFreebsd
  PlatformLinux
  PlatformNacl
  PlatformNetbsd
  PlatformOpenbsd
  PlatformPlan9
  PlatformSolaris
  PlatformWindows
)

func KPlatformLook() int {

  switch runtime.GOOS {
  case "android":
    return PlatformAndroid
  case "darwin":
    return PlatformDarwin
  case "dragonfly":
    return PlatformDragonfly
  case "freebsd":
    return PlatformFreebsd
  case "linux":
    return PlatformLinux
  case "nacl":
    return PlatformNacl
  case "netbsd":
    return PlatformNetbsd
  case "openbsd":
    return PlatformOpenbsd
  case "plan9":
    return PlatformPlan9
  case "solaris":
    return PlatformSolaris
  case "windows":
    return PlatformWindows
  }

  return PlatformUnknown
}
