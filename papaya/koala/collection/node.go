package collection

type KNode[T comparable] struct {
  prev *KNode[T]
  next *KNode[T]

  value T
}

type KNodeImpl[T comparable] interface {
  Init(value T)
  Set(value T)
  After(node *KNode[T])
  Before(node *KNode[T])
  Swap(node *KNode[T])
  Next() *KNode[T]
  Prev() *KNode[T]
  Value() T
  Free()
}

func KNodeNew[T comparable](value T) *KNode[T] {

  node := &KNode[T]{}
  node.Init(value)

  return node
}

func (v *KNode[T]) Init(value T) {

  v.prev = nil
  v.next = nil
  v.value = value
}

func (v *KNode[T]) Set(value T) {

  v.value = value
}

func (v *KNode[T]) After(node *KNode[T]) {

  var next *KNode[T]

  if node != nil {

    next = node.next

    // linked on prev
    v.prev = node
    node.next = v

    // linked on next
    v.next = next

    if next != nil {
      next.prev = v
    }
  }
}

func (v *KNode[T]) Before(node *KNode[T]) {

  var prev *KNode[T]

  if node != nil {

    prev = node.prev

    // linked on next
    v.next = node
    node.prev = v

    // linked on prev
    v.prev = prev

    if prev != nil {
      prev.next = v
    }
  }
}

func (v *KNode[T]) Swap(node *KNode[T]) {

  var next, prev *KNode[T]

  if node != nil {

    next = node.next
    prev = node.prev

    // linked on current node
    node.next = v.next
    node.prev = v.prev

    // swap on next
    if next != nil {
      v.next = next
    }

    // swap on prev
    if prev != nil {
      v.prev = prev
    }
  }
}

func (v *KNode[T]) Next() *KNode[T] {

  return v.next
}

func (v *KNode[T]) Prev() *KNode[T] {

  return v.prev
}

func (v *KNode[T]) Value() T {

  return v.value
}

func (v *KNode[T]) Free() {

  var next, prev *KNode[T]

  next = v.next
  prev = v.prev

  // free memory
  v.next = nil
  v.prev = nil

  // linked on next to prev
  if next != nil {
    next.prev = prev
  }

  // linked on prev to next
  if prev != nil {
    prev.next = next
  }
}
