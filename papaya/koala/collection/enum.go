package collection

// ---------------------------- Enum, Enums ----------------------------

// enumeration, no required pointer or references

type KEnum[K any, V any] struct {
  key   K
  value V
}

type KEnumImpl[K any, V any] interface {
  Key() K
  Value() V
  Tuple() (K, V)
}

func KEnumNew[K any, V any](key K, value V) KEnumImpl[K, V] {

  return &KEnum[K, V]{
    key:   key,
    value: value,
  }
}

func (v KEnum[K, V]) Key() K {

  return v.key
}

func (v KEnum[K, V]) Value() V {

  return v.value
}

func (v KEnum[K, V]) Tuple() (K, V) {

  return v.key, v.value
}
