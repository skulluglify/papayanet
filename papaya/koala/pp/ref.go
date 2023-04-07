package pp

// implementation `references` instead `pointer` methods
// funny as well, make it hard to use

type KRef[T any] struct {
  Value *T
}

type KRefImpl[T any] interface {
  Init(value *T)
  Get() T
  Set(value T)
  Ptr() *T
  Free()
}

func KRefNew[T any](value T) KRefImpl[T] {

  ref := &KRef[T]{}
  ref.Set(value)

  return ref
}

func KRefNewPtr[T any](value *T) KRefImpl[T] {

  ref := &KRef[T]{}
  ref.Init(value)

  return ref
}

func (ref *KRef[T]) Get() T {

  return *ref.Value
}

func (ref *KRef[T]) Set(value T) {

  ref.Value = &value
}

func (ref *KRef[T]) Ptr() *T {

  return ref.Value
}

func (ref *KRef[T]) Init(value *T) {

  ref.Value = value
}

func (ref *KRef[T]) Free() {

  ref.Value = nil
}
