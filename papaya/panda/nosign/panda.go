package nosign

import "math"

func Ceil(size float64) uint {

  ex := size - float64(uint(size))

  if ex > 0 {

    return uint(size + .5)
  }

  return uint(size)
}

func Floor(size float64) uint {

  ex := size - float64(uint(size))

  if ex > 0 {

    return uint(size - .5)
  }

  return uint(size)
}

func CeilHalf(size uint) uint {

  if size > 0 {

    if size < math.MaxUint {

      if size&1 == 1 {

        size++
      }

      //} else {

      // half of uint
      //return 9223372036854775808
    }

    return size / 2
  }

  return 0
}

func FloorHalf(size uint) uint {

  if size > 0 {

    if size&1 == 1 {

      size--
    }

    return size / 2
  }

  return 0
}

func Min(nums ...uint) uint {

  var num uint
  num = math.MaxUint

  for _, v := range nums {

    if v < num {

      num = v
    }
  }

  return num
}

func Max(nums ...uint) uint {

  var num uint
  num = 0 // 0 is minimum

  for _, v := range nums {

    if num < v {

      num = v
    }
  }

  return num
}

func Avg(nums ...uint) uint {

  var num, k uint
  num = 0
  k = 0

  for _, v := range nums {

    if k == 1 {

      num += v
      num /= 2
    }
    num = v
    k = 1
  }

  return num
}
