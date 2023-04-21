package mapping

type Keys []string

func (v Keys) Contain(t string) bool {

  for _, k := range v {

    if k == t {

      return true
    }
  }

  return false
}
