package mapping

import (
  "skfw/papaya/koala/collection"
  "skfw/papaya/koala/gen"
)

// typed KMap with a commonly built-in type

type KMap map[string]any
type KMapImpl interface {
  Iterable() gen.KMapIterationImpl
  Enums() []collection.KEnumImpl[string, any]
  Keys() Keys
  Values() []any
  Get(name string) any
  Set(name string, data any) bool
  Del(name string) bool
  Branch(name string) any
  Put(name string, data any) bool
  JSON() string
  Tree() KMapTreeImpl
}

type KMapTree KMap
type KMapTreeImpl interface {
  Iterable() gen.KMapIterationImpl
  Enums() []collection.KEnumImpl[string, any]
  Keys() Keys
  Get(name string) any
  Set(name string, data any) bool
  Del(name string) bool
  Branch(name string) any
  Put(name string, data any) bool
  JSON() string
  Inline() KMapImpl
}

func (m *KMap) Iterable() gen.KMapIterationImpl {

  return gen.KMapIterable(m)
}

func (m *KMap) Enums() []collection.KEnumImpl[string, any] {

  return KMapEnums(m)
}

func (m *KMap) Keys() Keys {

  return KMapKeys(m)
}

func (m *KMap) Values() []any {

  return KMapValues(m)
}

func (m *KMap) Get(name string) any {

  return KMapGetValue(name, m)
}

func (m *KMap) Set(name string, data any) bool {

  return KMapSetValue(name, data, m)
}

func (m *KMap) Del(name string) bool {

  return KMapDelValue(name, m)
}

func (m *KMap) Branch(name string) any {

  return KMapBranch(name, m)
}

func (m *KMap) Put(name string, data any) bool {

  return KMapPut(name, data, m)
}

func (m *KMap) JSON() string {

  return KMapEncodeJSON(m)
}

func (m *KMap) Tree() KMapTreeImpl {

  tree := KMapTree(*m)
  return &tree
}

func (m *KMapTree) Iterable() gen.KMapIterationImpl {

  return gen.KMapTreeIterable(m)
}

func (m *KMapTree) Enums() []collection.KEnumImpl[string, any] {

  return KMapTreeEnums(m)
}

func (m *KMapTree) Keys() Keys {

  return KMapTreeKeys(m)
}

func (m *KMapTree) Get(name string) any {

  return KMapGetValue(name, m)
}

func (m *KMapTree) Set(name string, data any) bool {

  return KMapSetValue(name, data, m)
}

func (m *KMapTree) Del(name string) bool {

  return KMapDelValue(name, m)
}

func (m *KMapTree) Branch(name string) any {

  return KMapBranch(name, m)
}

func (m *KMapTree) Put(name string, data any) bool {

  return KMapPut(name, data, m)
}

func (m *KMapTree) JSON() string {

  return KMapEncodeJSON(m)
}

func (m *KMapTree) Inline() KMapImpl {

  inline := KMap(*m)
  return &inline
}
