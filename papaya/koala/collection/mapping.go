package collection

// Layer List
// Weird projekt coming soon

type MapList[T comparable] struct {
  data KListImpl[KEnumImpl[string, T]]
}

type MapListImpl[T comparable] interface {
  Init()
}

func (l *MapList[T]) Init() {

  panic("not implemented yet")
}
