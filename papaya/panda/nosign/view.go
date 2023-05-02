package nosign

import "strconv"

func Look(a uint64, b uint64, prec int) string {

  n, d := float64(a), float64(b)
  k := n / d
  z := uint64(n / d)
  m := k - float64(z)
  c := uint64(m * 10)

  if c == 0 {

    return strconv.FormatUint(z, 10)
  }

  return strconv.FormatFloat(k, 'f', prec, 32)
}
