package collection

// Layer List
// Weird projekt coming soon

type LayerList[T comparable] struct {
  layers KListImpl[KListImpl[T]]
}

type LayerListImpl[T comparable] interface {
  Init()
}

func (l *LayerList[T]) Init() {

  panic("not implemented yet")
}
