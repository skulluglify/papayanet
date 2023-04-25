package collection

import (
	"PapayaNet/papaya/panda/nosign"
	"errors"
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
	removeNodeByIndex(index uint) error
	removeNode(node *KNode[T])

	refresh(pos int, index uint) // only increase
	reset()
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
		m.reset(MoveUp)
	}
}

func (m *KMiddleList[T]) Remove(values ...T) error {

	var n uint

	n = uint(len(values))

	if n > 0 {

		var err error
		var actionMove ActionMoveImpl

		// decrease
		if actionMove, err = m.uniMoveDown(0, n); err == nil {

			err = m.updateMiddlePosition(actionMove)
		}

		if err != nil {

			panic(err)
		}
	}

	return m.array.Remove(values...)
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

func (m *KMiddleList[T]) Splice(index uint, deleteCount uint, values ...T) KListImpl[T] {

	var err error
	var actionMove ActionMoveImpl

	// decrease
	if deleteCount > 0 {

		if actionMove, err = m.uniMoveDown(index, deleteCount); err == nil {

			err = m.updateMiddlePosition(actionMove)
		}

		if err != nil {

			panic(err)
		}
	}

	// splice
	var value T
	var i, j, n uint
	var array KListImpl[T]
	var nodeSelect, nodeSafe, node *KNode[T]

	array = KMiddleListNew[T]()
	nodeSelect, err = m.findNodeByIndex(index)

	if err != nil {

		array.Add(values...)
		return array
	}

	n = uint(len(values))

	nodeSafe = nodeSelect

	for i = 0; i < n; i++ {

		j = n - i - 1

		value = values[j]
		node = KNodeNew(value)

		// look like `PushLeft`
		node.Before(nodeSafe)

		if nodeSafe == m.array.head {

			m.array.head = node
		}

		nodeSafe = node
		m.array.size += 1
	}

	nodeSafe = nodeSelect

	for i = 0; i < deleteCount; i++ {

		if nodeSafe != nil {
			array.Push(nodeSafe.Value())
			node = nodeSafe.Next()

			// m.removeNode(nodeSafe)
			// node, k is determine

			// decrease, update middle pointing
			if actionMove, err = m.uniMoveDown(index, 1); err == nil {

				err = m.updateMiddlePosition(actionMove)
			}

			if err != nil {

				panic(err)
			}

			m.array.removeNode(nodeSafe) // without update middle pointing

			nodeSafe = node
			continue
		}
		break
	}

	//if n > 0 {
	//
	//  m.middle = m.array.size
	//  m.pos = 0
	//  m.refresh(0)
	//}

	return array
}

func (m *KMiddleList[T]) Copy() KListImpl[T] {

	var i uint
	var array KListImpl[T]
	var node KNodeImpl[T]

	array = KListNew[T]()

	node = m.array.head

	for i = 0; i < m.array.size; i++ {

		array.Push(node.Value())

		node = node.Next()
	}

	return array
}

func (m *KMiddleList[T]) PushLeft(value T) {

	m.array.PushLeft(value)
	m.refresh(MoveUp, 0)
}

func (m *KMiddleList[T]) Push(value T) {

	m.array.Push(value)
	m.refresh(MoveUp, m.array.size-1)
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

func (m *KMiddleList[T]) ForEach(cb KListMapHandler[T]) error {

	return m.array.ForEach(cb)
}

func (m *KMiddleList[T]) uniMoveUp(index uint, size uint) (ActionMoveImpl, error) {

	var currLen, currMidPos, virtLen, virtMidPos uint

	// current position
	currLen = m.array.size
	currMidPos = m.pos

	var err error

	// virtual position
	virtLen = currLen + size
	virtMidPos, err = VirtMidPos(virtLen)

	if err != nil {

		return nil, err
	}

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

	if size > 0 {

		if index <= currMidPos {

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
	}

	return ActionMoveNew(pos, ran), nil
}

func (m *KMiddleList[T]) uniMoveDown(index uint, size uint) (ActionMoveImpl, error) {

	var currLen, currMidPos, virtLen, virtMidPos uint

	if m.array.size < size {

		return nil, errors.New("down size is larger than array size")
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

	virtMidPos, _ = VirtMidPos(virtLen)

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

	if size > 0 {

		if index <= virtMidPos {

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
	}

	return ActionMoveNew(pos, ran), nil
}

func (m *KMiddleList[T]) updateMiddlePosition(actionMove ActionMoveImpl) error {

	var k uint

	for k = actionMove.Range(); 0 < k; k-- {

		// unfinished
		if m.middle == nil {

			return errors.New("middle is NULL")
		}

		switch actionMove.Pos() {
		case MoveLeft:

			m.middle = m.middle.Prev()
			m.pos--
			break

		case MoveRight:

			m.middle = m.middle.Next()
			m.pos++
			break
		}
	}

	return nil
}

func (m *KMiddleList[T]) findNodeByIndex(index uint) (*KNode[T], error) {

	var err error
	var i, k, q, t uint
	var node *KNode[T]

	if m.array.size == 0 {

		return nil, errors.New("empty list")
	}

	if m.array.size <= index {

		return nil, errors.New("index out of bound")
	}

	node = nil

	k, err = VirtMidPos(m.array.size)
	if err != nil {

		return nil, err
	}

	q, err = VirtMidPos(nosign.FloorHalf(m.array.size))
	if err != nil {

		return nil, err
	}

	// t, err = VirtMidPos(nosign.FloorHalf(m.array.size))
	// if err != nil {

	// 	return nil, err
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

		var err error
		var actionMove ActionMoveImpl

		// decrease
		if actionMove, err = m.uniMoveDown(index, 1); err == nil {

			err = m.updateMiddlePosition(actionMove)
		}

		if err != nil {

			panic(err)
		}

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

			m.refresh(MoveDown, k)

		} else { // update by knowing index

			m.reset(MoveDown)

		}

		// free-up node
		node.Free()

		// update size
		m.array.size--
	}
}

func (m *KMiddleList[T]) refresh(pos int, index uint) {

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

		if actionMove, err = m.uniMoveUp(index, 1); err == nil {

			err = m.updateMiddlePosition(actionMove)
		}

	case MoveDown:

		if actionMove, err = m.uniMoveDown(index, 1); err == nil {

			err = m.updateMiddlePosition(actionMove)
		}
	}

	if err != nil {

		panic(err)
	}
}

func (m *KMiddleList[T]) reset(pos int) {

	// hard reset value
	m.middle = m.array.head
	m.pos = 0

	m.refresh(pos, 0)
}
