package panda

import "math"

func Abs(value int) int {

  if value > 0 {

    return value
  }

  return -value
}

func Ceil(size float64) int {

  ex := size - float64(int(size))

  if ex > 0 {

    return int(size + .5)
  }

  return int(size)
}

func Floor(size float64) int {

  ex := size - float64(int(size))

  if ex > 0 {

    return int(size - .5)
  }

  return int(size)
}

func CeilHalf(size int) int {

  if size < math.MaxInt {

    if size&1 == 1 {

      size++
    }

    //} else {

    // half of int
    //return 2147483647
  }

  return size / 2
}

func FloorHalf(size int) int {

  if size&1 == 1 {

    size--
  }

  return size / 2
}

func Min(nums ...int) int {

  num := math.MaxInt

  for _, v := range nums {

    if v < num {

      num = v
    }
  }

  return num
}

func Max(nums ...int) int {

  num := math.MinInt

  for _, v := range nums {

    if num < v {

      num = v
    }
  }

  return num
}

func Avg(nums ...int) int {

  num := 0
  k := 0

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
