package collection

import (
  "errors"
  "skfw/papaya/panda/nosign"
)

type KMiddleList[T comparable] struct {
  array  *KList[T]
  middle *KNode[T]

  // middle position is relative
  // can't use it for next middle position
  pos uint
}

type KMiddleListImpl[T comparable] interface {
  KListImpl[T]

  // Universal Movement Increase, Decrease

  uniMoveUp(index uint, size uint) (ActionMoveImpl, error)
  uniMoveDown(index uint, size uint) (ActionMoveImpl, error)
  updateMiddlePosition(actionMove ActionMoveImpl) error

  // Utility

  findNodeByIndex(index uint) (*KNode[T], error)
  findNodeByValue(value T) (*KNode[T], error)

  removeNodeByIndex(index uint) error
  removeNode(node *KNode[T])

  refresh(pos int, index uint, size uint) // only increase
  reset(pos int, size uint)
}

func KMiddleListNew[T comparable]() KListImpl[T] {

  array := &KMiddleList[T]{}
  array.Init()

  return array
}

func KMiddleListNewR[T comparable](data []T) KListImpl[T] {

  array := KMiddleListNew[T]()

  var i, n uint

  n = uint(len(data))

  for i = 0; i < n; i++ {

    // copy value into an array array
    array.Add(data[i])
  }

  return array
}

func KMiddleListNewV[T comparable](values ...T) KListImpl[T] {

  array := KMiddleListNew[T]()

  // set variadic values
  array.Add(values...)

  return array
}

func (m *KMiddleList[T]) Init() {

  m.array = &KList[T]{}
  m.middle = nil
  m.pos = 0
}

func (m *KMiddleList[T]) Add(values ...T) {

  var n uint

  n = uint(len(values))

  if n > 0 {

    m.array.Add(values...)
    m.reset(MoveUp, 0)
  }
}

func (m *KMiddleList[T]) Remove(values ...T) error {

  var err error
  var i, n uint
  var value T
  var node *KNode[T]

  n = uint(len(values))

  for i = 0; i < n; i++ {

    value = values[i]

    if node, err = m.findNodeByValue(value); err != nil {

      return err
    }

    m.removeNode(node)
  }

  return nil
}

func (m *KMiddleList[T]) Includes(values ...T) bool {

  return m.array.Includes(values...)
}

func (m *KMiddleList[T]) Get(index uint) (T, error) {

  node, err := m.findNodeByIndex(index)
  if err != nil {

    var noop T
    return noop, err
  }

  return node.Value(), nil
}

func (m *KMiddleList[T]) Set(index uint, value T) error {

  node, err := m.findNodeByIndex(index)
  if err != nil {

    return err
  }

  node.Set(value)
  return nil
}

func (m *KMiddleList[T]) Del(index uint) error {

  var err error
  var actionMove ActionMoveImpl

  // decrease
  if actionMove, err = m.uniMoveDown(index, 1); err == nil {

    err = m.updateMiddlePosition(actionMove)
  }

  if err != nil {

    return err
  }

  var node *KNode[T]

  node, err = m.findNodeByIndex(index)
  if err != nil {

    return err
  }

  m.removeNode(node)

  return nil
}

func (m *KMiddleList[T]) Slice(index uint, size uint) (KListImpl[T], error) {

  var err error
  var array KListImpl[T]
  var node KNodeImpl[T]

  array = KMiddleListNew[T]()

  var i, n uint

  node, err = m.findNodeByIndex(index)

  if err != nil {

    return array, err
  }

  n = size + index

  if m.array.size < n {

    return array, errors.New("size out of range")
  }

  for i = index; i < n; i++ {

    array.Push(node.Value())
    node = node.Next()
  }

  return array, nil
}

func (m *KMiddleList[T]) Splice(index uint, deleteCount uint, values ...T) (KListImpl[T], error) {

  var err error

  var array KListImpl[T]

  m.refresh(MoveDown, index, deleteCount)

  array, err = m.array.Splice(index, deleteCount, values...)

  m.refresh(MoveUp, index, uint(len(values)))

  return array, err
}

func (m *KMiddleList[T]) Copy() KListImpl[T] {

  var i uint
  var array KListImpl[T]
  var node KNodeImpl[T]

  array = KMiddleListNew[T]()

  node = m.array.head

  for i = 0; i < m.array.size; i++ {

    array.Push(node.Value())

    node = node.Next()
  }

  return array
}

func (m *KMiddleList[T]) PushLeft(value T) {

  m.array.PushLeft(value)
  m.refresh(MoveUp, 0, 1)
}

func (m *KMiddleList[T]) Push(value T) {

  m.array.Push(value)
  m.refresh(MoveUp, m.array.size-1, 1)
}

func (m *KMiddleList[T]) PopLeft() (T, error) {

  if m.array.size <= 0 {

    var noop T
    return noop, errors.New("empty list")
  }

  var err error
  var actionMove ActionMoveImpl

  // decrease
  if actionMove, err = m.uniMoveDown(0, 1); err == nil {

    err = m.updateMiddlePosition(actionMove)
  }

  var value T
  value = m.array.head.Value()
  if err = m.array.Del(0); err != nil {

    var noop T
    return noop, err
  }

  return value, nil
}

func (m *KMiddleList[T]) Pop() (T, error) {

  if m.array.size <= 0 {

    var noop T
    return noop, errors.New("empty list")
  }

  var err error
  var actionMove ActionMoveImpl

  // decrease
  if actionMove, err = m.uniMoveDown(m.array.size-1, 1); err == nil {

    err = m.updateMiddlePosition(actionMove)
  }

  if err != nil {

    panic(err)
  }

  var value T
  value = m.array.tail.Value()
  if err = m.array.Del(m.array.size - 1); err != nil {

    var noop T
    return noop, err
  }

  return value, nil
}

func (m *KMiddleList[T]) Len() uint {

  return m.array.size
}

func (m *KMiddleList[T]) Concat(values ...T) KListImpl[T] {

  var array KListImpl[T]
  array = m.Copy()
  array.Add(values...)
  return array
}

func (m *KMiddleList[T]) ConcatArray(arrays ...KListImpl[T]) (KListImpl[T], error) {

  var err error

  var i uint
  var array KListImpl[T]
  var value T
  array = m.Copy()

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

func (m *KMiddleList[T]) Replace(array KListImpl[T]) error {

  var err error

  var i uint
  var node *KNode[T]
  var value T

  node = m.array.head

  for i = 0; i < nosign.Min(m.array.size, array.Len()); i++ {

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
    m.Push(value)
  }

  return nil
}

func (m *KMiddleList[T]) Reverse() {

  // reverse
  m.array.Reverse()

  // update middle position
  if m.array.size&1 == 0 {
    m.pos += 1
  }

  // update middle node position
  m.refresh(MoveUp, 0, 0) // refresh middle
}

func (m *KMiddleList[T]) ForEach(cb KListMapHandler[T]) error {

  return m.array.ForEach(cb)
}

func (m *KMiddleList[T]) uniMoveUp(index uint, size uint) (ActionMoveImpl, error) {

  var currLen, currMidPos, virtLen, virtMidPos uint

  // current position
  currLen = m.array.size
  currMidPos = m.pos

  // virtual position
  virtLen = currLen + size

  if virtLen == 0 {

    return nil, errors.New("array is empty")
  }

  virtMidPos = VirtMidPos(virtLen)

  // determine the value of the position and distance
  var pos int
  var ran uint

  // default value
  pos = MoveLeft
  ran = 0

  // the next middle position is a relative
  // if the action takes increase before a middle
  // the next middle shifted by left

  // post - processing

  // range current middle into the next middle
  // [ O O O X X X O O O O O O O ] 10 -> 13
  //                 ^ currMidPos 4 -> 7
  //               ^ virtMidPos 6 -> left 1

  if index <= currMidPos {

    // pre - processing
    currMidPos = currMidPos + size // skipping
  }

  // relative increase or decrease
  if virtMidPos < currMidPos {

    pos = MoveLeft
    ran = currMidPos - virtMidPos

  } else {

    pos = MoveRight
    ran = virtMidPos - currMidPos
  }

  return ActionMoveNew(pos, ran), nil
}

func (m *KMiddleList[T]) uniMoveDown(index uint, size uint) (ActionMoveImpl, error) {

  var currLen, currMidPos, virtLen, virtMidPos uint

  if m.array.size < size {

    return nil, errors.New("downsize is larger than array size")
  }

  // current position
  currLen = m.array.size
  currMidPos = m.pos

  // virtual position
  virtLen = currLen - size

  if virtLen == 0 { // downside come nullable

    // pre - processing
    m.middle = nil
    m.pos = 0

    return nil, nil
  }

  virtMidPos = VirtMidPos(virtLen)

  // determine the value of the position and distance
  var pos int
  var ran uint

  // default value
  pos = MoveLeft
  ran = 0

  // the next middle position is a relative
  // if the action takes increase before a middle
  // the next middle shifted by right

  // pre - processing

  // range current middle into the next middle
  // [ O O O X X X O O O O O O O ] 13 -> 10
  //               ^ currMidPos 6
  //                 ^ virtMidPos 4 -> 7 // real of virt middle position

  if index <= virtMidPos {

    // post - processing
    virtMidPos = virtMidPos + size // skipping
  }

  // relative increase or decrease
  if virtMidPos < currMidPos {

    pos = MoveLeft
    ran = currMidPos - virtMidPos

  } else {

    pos = MoveRight
    ran = virtMidPos - currMidPos
  }

  return ActionMoveNew(pos, ran), nil
}

func (m *KMiddleList[T]) updateMiddlePosition(actionMove ActionMoveImpl) error {

  var k uint

  // unfinished
  if m.middle == nil {

    return errors.New("middle is NULL")
  }

  switch actionMove.Position() {
  case MoveLeft:

    for k = 0; k < actionMove.Range(); k++ {

      m.middle = m.middle.Prev()
      m.pos--

    }
    break

  case MoveRight:

    for k = 0; k < actionMove.Range(); k++ {

      m.middle = m.middle.Next()
      m.pos++

    }
    break
  }

  return nil
}

func (m *KMiddleList[T]) findNodeByIndex(index uint) (*KNode[T], error) {

  var i, k, q, t uint
  var node *KNode[T]

  if m.array.size == 0 {

    return nil, errors.New("empty list")
  }

  if m.array.size <= index {

    return nil, errors.New("index out of bound")
  }

  node = nil

  //k = VirtMidPos(m.array.size)
  k = m.pos
  //q = VirtMidPos(nosign.FloorHalf(m.array.size))
  q = VirtMidPos(k + 1)

  // t, err = VirtMidPos(nosign.FloorHalf(m.array.size))
  // if err != nil {

  //   return nil, err
  // }

  // t = t + k
  t = k + q

  if index <= q {

    node = m.array.head

    for i = 0; i < index; i++ {

      node = node.Next()
    }

  } else { // if q < index

    if index <= k {

      node = m.middle

      for i = k - 1; index <= i; i-- {

        node = node.Prev()
      }

    } else { // if k < index

      if index <= t {

        node = m.middle

        for i = k; i < index; i++ {

          node = node.Next()
        }

      } else { // if t < index

        node = m.array.tail

        for i = m.array.size - 1; index < i; i-- {

          node = node.Prev()
        }
      }
    }
  }

  return node, nil
}

func (m *KMiddleList[T]) findNodeByValue(value T) (*KNode[T], error) {

  return m.array.findNodeByValue(value)
}

func (m *KMiddleList[T]) removeNodeByIndex(index uint) error {

  node, err := m.findNodeByIndex(index)
  if err != nil {

    return err
  }

  // var `k` has been determine by index

  // before free-up node
  // prevent head, middle, and tail
  // safety node removed by size
  if m.array.size > 0 {

    if node == m.array.head {

      m.array.head = node.Next()
    }

    if node == m.array.tail {

      m.array.tail = node.Prev()
    }

    m.refresh(MoveDown, index, 1)

    // free-up node
    node.Free()

    // update size
    m.array.size--
  }

  return nil
}

func (m *KMiddleList[T]) removeNode(node *KNode[T]) {

  var k uint
  var found bool
  // before free-up node
  // prevent head, middle, and tail
  // safety node removed by size
  if m.array.size > 0 {

    if node == m.array.head {

      m.array.head = node.Next()
      found = true
      k = 0
    }

    if node == m.middle {

      found = true
      k = m.pos
    }

    if node == m.array.tail {

      m.array.tail = node.Prev()
      k = m.array.size - 1
      found = true
    }

    if found { // hard reset middle position

      m.refresh(MoveDown, k, 1)

    } else { // update by knowing index

      m.reset(MoveDown, 1) // index?

    }

    // free-up node
    node.Free()

    // update size
    m.array.size--
  }
}

func (m *KMiddleList[T]) refresh(pos int, index uint, size uint) {

  if m.middle == nil {

    m.middle = m.array.head
    m.pos = 0
  }

  var err error
  var actionMove ActionMoveImpl

  // default value
  err = nil

  switch pos {
  case MoveUp:

    if actionMove, err = m.uniMoveUp(index, size); err == nil {

      err = m.updateMiddlePosition(actionMove)
    }

  case MoveDown:

    if actionMove, err = m.uniMoveDown(index, size); err == nil {

      err = m.updateMiddlePosition(actionMove)
    }
  }

  if err != nil {

    panic(err)
  }
}

func (m *KMiddleList[T]) reset(pos int, size uint) {

  // hard reset value
  m.middle = m.array.head
  m.pos = 0

  m.refresh(pos, 0, size)
}
