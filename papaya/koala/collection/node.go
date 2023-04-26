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
  Gt(node *KNode[T]) bool
  Ge(node *KNode[T]) bool
  Lt(node *KNode[T]) bool
  Le(node *KNode[T]) bool
  Eq(node *KNode[T]) bool
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

  // universal movement
  var next, prev *KNode[T]

  if node != nil {

    /*
    * Normalize
     */

    next = node.next // nullable
    prev = node.prev // nullable

    if next != nil {

      next.prev = v
    }

    if prev != nil {

      prev.next = v
    }

    /*
    * Normalize
     */

    // linked on current node
    node.next = v.next // nullable
    node.prev = v.prev // nullable

    // swap on another node
    v.next = next // nullable
    v.prev = prev // nullable

    /*
    * Normalize
     */

    next = node.next // is v.next
    prev = node.prev // is v.prev

    if next != nil {

      next.prev = node
    }

    if prev != nil {

      prev.next = node
    }

    /*
    * Normalize
     */
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

// comparable

func (v *KNode[T]) Gt(node *KNode[T]) bool {

  return CompareNew[T](v.value, node.value).Gt()
}

func (v *KNode[T]) Ge(node *KNode[T]) bool {

  return CompareNew[T](v.value, node.value).Ge()
}

func (v *KNode[T]) Lt(node *KNode[T]) bool {

  return CompareNew[T](v.value, node.value).Lt()
}

func (v *KNode[T]) Le(node *KNode[T]) bool {

  return CompareNew[T](v.value, node.value).Le()
}

func (v *KNode[T]) Eq(node *KNode[T]) bool {

  return CompareNew[T](v.value, node.value).Eq()
}
