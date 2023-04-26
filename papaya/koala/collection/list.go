package collection

import (
  "errors"
  "skfw/papaya/panda/nosign"
)

type KList[T comparable] struct {
  head *KNode[T]
  tail *KNode[T]

  size uint
}

type KListImpl[T comparable] interface {
  Init()
  Add(values ...T)
  Remove(values ...T) error
  Includes(values ...T) bool
  Get(index uint) (T, error)
  Set(index uint, value T) error
  Del(index uint) error
  Slice(index uint, size uint) (KListImpl[T], error)
  Splice(index uint, deleteCount uint, values ...T) (KListImpl[T], error)
  Copy() KListImpl[T]
  PushLeft(value T)
  Push(value T)
  PopLeft() (T, error)
  Pop() (T, error)
  Len() uint

  Concat(values ...T) KListImpl[T]
  ConcatArray(arrays ...KListImpl[T]) (KListImpl[T], error)
  Replace(array KListImpl[T]) error

  Reverse()

  ForEach(cb KListMapHandler[T]) error

  // Merge, Replace

  // Helper Methods
  findNodeByIndex(index uint) (*KNode[T], error)
  findNodeByValue(value T) (*KNode[T], error)

  removeNodeByIndex(index uint) error
  removeNode(node *KNode[T])
}

func KListNew[T comparable]() KListImpl[T] {

  array := &KList[T]{}
  array.Init()

  return array
}

func KListNewR[T comparable](data []T) KListImpl[T] {

  array := KListNew[T]()

  var i, n uint

  n = uint(len(data))

  for i = 0; i < n; i++ {

    // copy value into an array array
    array.Add(data[i])
  }

  return array
}

func KListNewV[T comparable](values ...T) KListImpl[T] {

  array := KListNew[T]()

  // set variadic values
  array.Add(values...)

  return array
}

func (v *KList[T]) Init() {

  // initial default value
  v.head = nil
  v.tail = nil
  v.size = 0
}

func (v *KList[T]) Add(values ...T) {

  var size uint
  var start, end uint
  var nodeStart, nodeEnd *KNode[T]

  size = uint(len(values))

  if size > 0 {

    start = 0
    end = size - 1

    nodeStart = KNodeNew[T](values[start])
    nodeEnd = nodeStart

    if v.tail != nil {

      nodeStart.After(v.tail)

    } else {

      v.head = nodeStart
    }

    for i := start + 1; i < end; i++ {

      node := KNodeNew[T](values[i])
      nodeEnd = node

      // linked
      node.After(nodeStart)
      nodeStart = node
    }

    if end > 0 {

      nodeEnd = KNodeNew[T](values[end])

      // linked
      nodeEnd.After(nodeStart)
    }

    v.tail = nodeEnd
  }

  v.size += size
}

func (v *KList[T]) Remove(values ...T) error {

  var err error
  var i, n uint
  var value T
  var node *KNode[T]

  n = uint(len(values))

  for i = 0; i < n; i++ {

    value = values[i]

    if node, err = v.findNodeByValue(value); err != nil {

      return err
    }

    v.removeNode(node)
  }

  return nil
}

func (v *KList[T]) Includes(values ...T) bool {

  var value T
  var i, j, size uint
  var node KNodeImpl[T]
  var found bool

  if v.head != nil {

    size = uint(len(values))

    for i = 0; i < size; i++ {

      node = v.head
      value = values[i]
      found = false

      if node != nil {

        for j = 0; j < v.size; j++ {

          // need comparable value
          if node.Value() == value {

            found = true
            break
          }

          node = node.Next()
        }

        if !found {

          // node is not found
          return false
        }

        // has been found
        continue
      }

      // node is null
      return false
    }

    // all includes
    return true
  }

  // empty list
  return false
}

func (v *KList[T]) Get(index uint) (T, error) {

  node, err := v.findNodeByIndex(index)

  if err != nil {

    var noop T // get zero value
    return noop, err
  }

  return node.Value(), nil
}

func (v *KList[T]) Set(index uint, value T) error {

  node, err := v.findNodeByIndex(index)

  if err != nil {

    return err
  }

  node.Set(value)

  return nil
}

func (v *KList[T]) Del(index uint) error {

  node, err := v.findNodeByIndex(index)

  if err != nil {

    return err
  }

  v.removeNode(node)

  return nil
}

func (v *KList[T]) Slice(index uint, size uint) (KListImpl[T], error) {

  var err error
  var array KListImpl[T]
  var node KNodeImpl[T]

  array = KListNew[T]()

  var i, n uint

  node, err = v.findNodeByIndex(index)

  if err != nil {

    return array, err
  }

  n = size + index

  if v.size < n {

    return array, errors.New("size out of range")
  }

  for i = index; i < n; i++ {

    array.Push(node.Value())
    node = node.Next()
  }

  return array, nil
}

func (v *KList[T]) Splice(index uint, deleteCount uint, values ...T) (KListImpl[T], error) {

  var err error

  var i, j, n uint
  var array KListImpl[T]
  var nodeSelect, nodeSafe, node *KNode[T]
  var value T

  array = KListNew[T]()
  nodeSelect, err = v.findNodeByIndex(index)

  if err != nil {

    v.Add(values...)
    return array, err
  }

  n = uint(len(values))

  nodeSafe = nodeSelect

  var isHead bool

  if nodeSafe == v.head {

    isHead = true
  }

  for i = 0; i < n; i++ {

    j = n - i - 1

    value = values[j]
    node = KNodeNew(value)

    // look like `PushLeft`
    node.Before(nodeSafe)

    // set v.head if possible
    if isHead {

      v.head = node
    }

    if nodeSafe == v.head {

      v.head = node
    }

    nodeSafe = node
    v.size += 1
  }

  nodeSafe = nodeSelect

  for i = 0; i < deleteCount; i++ {

    if nodeSafe != nil {
      array.Push(nodeSafe.Value())
      node = nodeSafe.Next()
      v.removeNode(nodeSafe)
      nodeSafe = node
      continue
    }
    break
  }

  return array, nil
}

func (v *KList[T]) Copy() KListImpl[T] {

  var i uint
  var array KListImpl[T]
  var node KNodeImpl[T]

  array = KListNew[T]()

  node = v.head

  for i = 0; i < v.size; i++ {

    array.Push(node.Value())

    node = node.Next()
  }

  return array
}

func (v *KList[T]) Push(value T) {

  var node *KNode[T]

  node = KNodeNew[T](value)

  if v.tail != nil {

    node.After(v.tail)
    v.tail = node

  } else {

    v.tail = node
    v.head = v.tail
    v.size = 0
  }

  v.size += 1
}

func (v *KList[T]) PushLeft(value T) {

  var node *KNode[T]

  node = KNodeNew[T](value)

  if v.head != nil {

    node.Before(v.head)
    v.head = node

  } else {

    v.head = node
    v.tail = v.head
    v.size = 0
  }

  v.size += 1
}

func (v *KList[T]) Pop() (T, error) {

  var value T
  var node KNodeImpl[T]

  if v.tail == nil {

    return value, errors.New("empty list")
  }

  node = v.tail
  value = node.Value()

  // ---------- Tail ----------
  // absolute step to execute

  v.tail = node.Prev()

  node.Free()
  // ---------- Tail ----------

  v.size--

  return value, nil
}

func (v *KList[T]) PopLeft() (T, error) {

  var value T
  var node KNodeImpl[T]

  if v.head == nil {

    return value, errors.New("empty list")
  }

  node = v.head
  value = node.Value()

  // ---------- Head ----------
  // absolute step to execute

  v.head = node.Next()

  node.Free()
  // ---------- Head ----------

  v.size--

  return value, nil
}

func (v *KList[T]) Len() uint {

  return v.size
}

func (v *KList[T]) Concat(values ...T) KListImpl[T] {

  var array KListImpl[T]
  array = v.Copy()

  for _, value := range values {

    array.Push(value)
  }

  return array
}

func (v *KList[T]) ConcatArray(arrays ...KListImpl[T]) (KListImpl[T], error) {

  var err error

  var i uint
  var array KListImpl[T]
  var value T
  array = v.Copy()

  for _, arr := range arrays {

    for i = 0; i < arr.Len(); i++ {

      if value, err = arr.Get(i); err != nil {

        return array, err
      }

      array.Push(value)
    }
  }

  return array, nil
}

func (v *KList[T]) Replace(array KListImpl[T]) error {

  var err error

  var i uint
  var node *KNode[T]
  var value T

  node = v.head

  for i = 0; i < nosign.Min(v.size, array.Len()); i++ {

    // getting current value
    if value, err = array.Get(i); err != nil {

      return err
    }

    // set current node
    node.Set(value)

    // update node
    node = node.next
  }

  for i = i + 1; i < array.Len(); i++ {

    // getting current value
    if value, err = array.Get(i); err != nil {

      return err
    }

    // added current value
    v.Push(value)
  }

  return nil
}

func (v *KList[T]) findNodeByIndex(index uint) (*KNode[T], error) {

  var i, m uint
  var node *KNode[T]

  if v.size == 0 {

    return nil, errors.New("empty list")
  }

  if v.size <= index {

    return nil, errors.New("index out of bound")
  }

  m = nosign.CeilHalf(v.size)

  if index <= m {

    node = v.head

    for i = 1; i <= index; i++ {

      node = node.Next()
    }

  } else {

    node = v.tail

    //for i = v.size - 2; index <= i; i-- {
    //
    //  node = node.Prev()
    //}

    // safety
    for i = index; i+1 < v.size; i++ {

      node = node.Prev()
    }
  }

  return node, nil
}

func (v *KList[T]) findNodeByValue(value T) (*KNode[T], error) {

  var node *KNode[T]

  if v.size == 0 {

    return nil, errors.New("empty list")
  }

  node = v.head

  for node != v.tail {

    if node.Value() == value {

      return node, nil
    }

    node = node.Next()
  }

  return nil, errors.New("value has been not found")
}

func (v *KList[T]) removeNodeByIndex(index uint) error {

  node, err := v.findNodeByIndex(index)
  if err != nil {

    return err
  }

  v.removeNode(node)

  return nil
}

func (v *KList[T]) removeNode(node *KNode[T]) {

  // ---------- Head, Tail ----------

  // before free-up node
  // prevent head, and tail
  // safety node removed by size
  if v.size > 0 {

    if node == v.head {

      v.head = node.Next()
    }

    if node == v.tail {

      v.tail = node.Prev()
    }

    // free-up node
    node.Free()

    // update size
    v.size--
  }

  // ---------- Head, Tail ----------
}

func (v *KList[T]) Reverse() {

  var node, nodeSwap, nodeTemp *KNode[T]

  node = v.head

  v.head = v.tail
  v.tail = node

  nodeTemp = node.next

  for nodeTemp != nil {

    // initial node
    node = nodeTemp.prev
    nodeSwap = nodeTemp

    // update nodeTemp
    nodeTemp = nodeTemp.next

    // swap
    node.prev = nodeSwap
    nodeSwap.next = node
  }

  // safe null
  v.head.prev = nil
  v.tail.next = nil
}

func (v *KList[T]) ForEach(cb KListMapHandler[T]) error {

  var i uint
  var next *KNode[T]
  i = 0

  next = v.head

  for next != nil {

    if err := cb(i, next.Value()); err != nil {

      return err
    }

    next = next.Next()
    i++
  }

  return nil
}
