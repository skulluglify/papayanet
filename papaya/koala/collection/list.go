package collection

import (
  "PapayaNet/papaya/panda/nosign"
  "errors"
  "fmt"
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
  Splice(index uint, deleteCount uint, values ...T) KListImpl[T]
  Copy() KListImpl[T]
  PushLeft(value T)
  Push(value T)
  PopLeft() (T, error)
  Pop() (T, error)
  Len() uint

  // Concat(values ...any) KListImpl[T] << KListImpl[T] / T
  // Merge, Replace

  // Helper Methods
  findNodeByIndex(index uint) (*KNode[T], error)
  removeNode(node *KNode[T])
}

func KListNew[T comparable]() KListImpl[T] {

  list := &KList[T]{}
  list.Init()

  return list
}

func KListNewR[T comparable](data []T) KListImpl[T] {

  list := &KList[T]{}
  list.Init()

  var i, n uint

  n = uint(len(data))

  for i = 0; i < n; i++ {

    // copy value into an array list
    list.Add(data[i])
  }

  return list
}

func KListNewV[T comparable](values ...T) KListImpl[T] {

  list := &KList[T]{}
  list.Init()

  // set variadic values
  list.Add(values...)

  return list
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

  var value T
  var i, j, size uint
  var node *KNode[T]
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

            v.removeNode(node)
            found = true
            break
          }

          node = node.Next()
        }

        if !found {

          return errors.New(fmt.Sprint("value", value, "is not found"))
        }

        // value has been removed
        continue
      }

      return errors.New("node is null")
    }

    // all values have been removed
    return nil
  }

  return errors.New("empty list")
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

  var i, j, m, n uint
  var list KListImpl[T]
  var node KNodeImpl[T]

  list = KListNew[T]()

  if v.size <= index {

    return list, errors.New("index out of bound")
  }

  n = size + index

  if v.size < n {

    return list, errors.New("size out of range")
  }

  m = nosign.CeilHalf(v.size)

  if index <= m {

    node = v.head

    for i = 0; i < n; i++ {

      if index <= i {

        list.Push(node.Value())
      }

      node = node.Next()
    }

  } else {

    node = v.tail

    //for i = v.size - 1; index <= i; i-- {
    //
    //  if i < n {
    //
    //    list.PushLeft(node.Value())
    //  }
    //
    //  node = node.Prev()
    //}

    // safety
    for i = index; i < v.size; i++ {

      j = v.size - i - 1

      if j < n {

        list.PushLeft(node.Value())
      }

      node = node.Prev()
    }
  }

  return list, nil
}

func (v *KList[T]) Splice(index uint, deleteCount uint, values ...T) KListImpl[T] {

  var value T
  var i, j, n uint
  var list KListImpl[T]
  var nodeSelect, nodeSafe, node *KNode[T]
  var err error

  list = KListNew[T]()
  nodeSelect, err = v.findNodeByIndex(index)

  if err != nil {

    v.Add(values...)
    return list
  }

  n = uint(len(values))

  nodeSafe = nodeSelect

  for i = 0; i < n; i++ {

    j = n - i - 1

    value = values[j]
    node = KNodeNew(value)

    // look like `PushLeft`
    node.Before(nodeSafe)

    if nodeSafe == v.head {

      v.head = node
    }

    nodeSafe = node
    v.size += 1
  }

  nodeSafe = nodeSelect

  for i = 0; i < deleteCount; i++ {

    if nodeSafe != nil {
      list.Push(nodeSafe.Value())
      node = nodeSafe.Next()
      v.removeNode(nodeSafe)
      nodeSafe = node
      continue
    }
    break
  }

  return list
}

func (v *KList[T]) Copy() KListImpl[T] {

  var i uint
  var list KListImpl[T]
  var node KNodeImpl[T]

  list = KListNew[T]()

  node = v.head

  for i = 0; i < v.size; i++ {

    list.Push(node.Value())

    node = node.Next()
  }

  return list
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

func (v *KList[T]) findNodeByIndex(index uint) (*KNode[T], error) {

  var i, m uint
  var node *KNode[T]

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
