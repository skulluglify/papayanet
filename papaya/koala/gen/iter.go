package gen

import (
  "PapayaNet/papaya/koala"
  "errors"
)

// ---------------------------- Iteration ----------------------------

type KIterationNextHandler[K any, V any] func(iter KIterationImpl[K, V]) error

type KIteration[K any, V any] struct {

  // enum
  key   K
  value V

  // next handling
  NextHandler KIterationNextHandler[K, V]

  // catch error from next handling
  err error

  // stop next iteration
  done bool

  // force stop iteration
  stopIter bool

  // marked if method `Set` has been called
  mark bool

  /*
     done stopIter
     0    0
     1    0
     1    1

  */
}

type KIterationImpl[K any, V any] interface {
  Set(key K, value V, done bool)
  SetValues(key K, value V)
  SetEnum(enum koala.KEnumImpl[K, V])
  Enum() koala.KEnumImpl[K, V]
  Next() KIterationImpl[K, V]
  HasNext() bool
  Error() error
  Stop() error
}

func (v *KIteration[K, V]) Set(key K, value V, done bool) {

  v.key = key
  v.value = value

  // once set if true
  // give signal `stopIter`
  if done {

    v.done = true
  }

  v.mark = true
}

func (v *KIteration[K, V]) SetValues(key K, value V) {

  v.key = key
  v.value = value

  // done already false
  // if var `done` is true, give signal `stopIter`
  //v.done = false

  v.mark = true
}

func (v *KIteration[K, V]) SetEnum(enum koala.KEnumImpl[K, V]) {

  v.key = enum.Key()
  v.value = enum.Value()

  // done already false
  // if var `done` is true, give signal `stopIter`
  //v.done = false

  v.mark = true
}

func (v *KIteration[K, V]) Enum() koala.KEnumImpl[K, V] {

  return koala.KEnumNew(v.key, v.value)
}

func (v *KIteration[K, V]) Error() error {

  return v.err
}

func (v *KIteration[K, V]) Stop() error {

  v.stopIter = true
  return errors.New("force stop iteration")
}

func (v *KIteration[K, V]) Next() KIterationImpl[K, V] {

  v.mark = false

  // canceled `stopIter` from signal done
  // if signal done set `false`
  v.stopIter = v.done

  if err := v.NextHandler(v); err != nil {

    v.err = err
    v.stopIter = true
  }

  // if method `Set` not called, force stop iteration
  if !v.mark {

    v.stopIter = true
  }

  return v
}

func (v *KIteration[K, V]) HasNext() bool {

  if v.stopIter {

    return false
  }

  return true
}

func KStopIter[K any, V any]() KIterationImpl[K, V] {

  return &KIteration[K, V]{
    stopIter: true,
  }
}
