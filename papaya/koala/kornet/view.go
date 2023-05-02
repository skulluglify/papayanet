package kornet

import (
  "skfw/papaya/panda/nosign"
  "strconv"
)

const EXABYTES uint64 = 1024 * 1024 * 1024 * 1024 * 1024 * 1024
const PETABYTES uint64 = 1024 * 1024 * 1024 * 1024 * 1024
const TERABYTES uint64 = 1024 * 1024 * 1024 * 1024
const GIGABYTES uint64 = 1024 * 1024 * 1024
const MEGABYTES uint64 = 1024 * 1024
const KILOBYTES uint64 = 1024

func ReprByte(size uint64) string {

  switch {

  case EXABYTES <= size:

    return nosign.Look(size, EXABYTES, 1) + "EB"

  case PETABYTES <= size:

    return nosign.Look(size, PETABYTES, 1) + "PB"

  case TERABYTES <= size:

    return nosign.Look(size, TERABYTES, 1) + "TB"

  case GIGABYTES <= size:

    return nosign.Look(size, GIGABYTES, 1) + "GB"

  case MEGABYTES <= size:

    return nosign.Look(size, MEGABYTES, 1) + "MB"

  case KILOBYTES <= size:

    return nosign.Look(size, KILOBYTES, 1) + "KB"

  }

  return strconv.FormatUint(size, 10) + "B"
}
