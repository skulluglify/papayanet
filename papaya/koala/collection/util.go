package collection

import "errors"

func VirtMidPos(size uint) (uint, error) {

  // size 0 = -1 is NULL
  if size == 0 {

    return 0, errors.New("empty list")
  }

  // size 1 = 0 | 2 = 0 1
  // size 2 = 0 | 2 = 0 1
  // size 3 = 1 | 4 = 1 2
  // size 4 = 1 | 4 = 1 2
  // size 5 = 2 | 6 = 2 3
  // size 6 = 2 | 6 = 2 3
  // size 7 = 3 | 8 = 3 4

  if size&1 == 1 {

    size += 1
  }

  size = size / 2
  size = size - 1

  return size, nil
}
