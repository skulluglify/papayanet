package runtime

import "runtime"

//goland:noinspection ALL
const (
  ArchUnknown = iota
  Arch386
  ArchAmd64
  ArchAmd64P32
  ArchArm
  ArchArmbe
  ArchArm64
  ArchArm64BE
  ArchPpc64
  ArchPpc64LE
  ArchMips
  ArchMipsle
  ArchMips64
  ArchMips64LE
  ArchMips64P32
  ArchMips64P32LE
  ArchPpc
  ArchRiscv64
  ArchS390X
)

func KArchLook() int {

  switch runtime.GOARCH {
  case "386":
    return Arch386
  case "amd64":
    return ArchAmd64
  case "amd64p32":
    return ArchAmd64P32
  case "arm":
    return ArchArm
  case "armbe":
    return ArchArmbe
  case "arm64":
    return ArchArm64
  case "arm64be":
    return ArchArm64BE
  case "ppc64":
    return ArchPpc64
  case "ppc64le":
    return ArchPpc64LE
  case "mips":
    return ArchMips
  case "mipsle":
    return ArchMipsle
  case "mips64":
    return ArchMips64
  case "mips64le":
    return ArchMips64LE
  case "mips64p32":
    return ArchMips64P32
  case "mips64p32le":
    return ArchMips64P32LE
  case "ppc":
    return ArchPpc
  case "riscv64":
    return ArchRiscv64
  case "s390x":
    return ArchS390X
  }

  return ArchUnknown
}
