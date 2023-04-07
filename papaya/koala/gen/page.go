package gen

type KPageIteration[K any, V any] struct {
  Page     []KIterationImpl[K, V]
  PageSize int
}

type KPageIterationImpl[K any, V any] interface {
  Init()
  Keys() []K
  Values() []V
  Add(iter KIterationImpl[K, V])
  Pop() KIterationImpl[K, V] // explicit method, remove end iteration list, and skipping next on previous iteration
  End() KIterationImpl[K, V]
  Wait() KIterationImpl[K, V]
  Len() int
}

func (v *KPageIteration[K, V]) Init() {

  v.Page = make([]KIterationImpl[K, V], 0)
  v.PageSize = 0
}

func (v *KPageIteration[K, V]) Keys() []K {

  keys := make([]K, 0)

  if v.PageSize > 0 {

    for i := 0; i < v.PageSize; i++ {

      iter := v.Page[i]

      if iter.HasNext() {

        enum := iter.Enum()
        key := enum.Key()
        keys = append(keys, key)
      }
    }
  }

  return keys
}

func (v *KPageIteration[K, V]) Values() []V {

  values := make([]V, 0)

  if v.PageSize > 0 {

    for i := 0; i < v.PageSize; i++ {

      iter := v.Page[i]

      if iter.HasNext() {

        enum := iter.Enum()
        value := enum.Value()
        values = append(values, value)
      }
    }
  }

  return values
}

func (v *KPageIteration[K, V]) Add(iter KIterationImpl[K, V]) {

  v.Page = append(v.Page, iter)
  v.PageSize += 1
}

func (v *KPageIteration[K, V]) Pop() KIterationImpl[K, V] {

  // slicing a pop iteration list
  if v.PageSize > 0 {

    v.Page = v.Page[:v.PageSize-1]
    v.PageSize += -1

    // skipping previous next iteration
    if v.PageSize > 0 {

      iter := v.Page[v.PageSize-1]
      //iter.Freeze()

      if iter.HasNext() {

        // next iteration
        next := iter.Next()
        v.Page[v.PageSize-1] = next

        return next
      }
    }
  }

  return KStopIter[K, V]()
}

func (v *KPageIteration[K, V]) End() KIterationImpl[K, V] {

  // catch the next iteration
  if v.PageSize > 0 {

    iter := v.Page[v.PageSize-1]
    //iter.Freeze()

    if iter.HasNext() {

      // next iteration
      next := iter.Next()
      v.Page[v.PageSize-1] = next

      return next
    }
  }

  return KStopIter[K, V]()
}

func (v *KPageIteration[K, V]) Wait() KIterationImpl[K, V] {

  next := v.End()

  for {

    if !next.HasNext() {

      if v.PageSize > 0 {

        next = v.Pop()
        continue
      }
    }

    break
  }

  return next
}

func (v *KPageIteration[K, V]) Len() int {

  return v.PageSize
}
