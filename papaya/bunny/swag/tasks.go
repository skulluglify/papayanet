package swag

import (
  "errors"
  "fmt"
  "skfw/papaya/koala/collection"
  m "skfw/papaya/koala/mapping"
)

type SwagTask struct {
  Name    string
  Handler SwagRouteHandler
}

func MakeSwagTask(name string, handler SwagRouteHandler) *SwagTask {

  return &SwagTask{
    Name:    name,
    Handler: handler,
  }
}

type SwagTasks collection.KListImpl[*SwagTask]

func MakeSwagTasks() collection.KListImpl[*SwagTask] {

  return collection.KListNew[*SwagTask]()
}

type SwagTasksQueue struct {
  tasks SwagTasks
}

type SwagTasksQueueImpl interface {
  Init()
  AddTask(task *SwagTask)
  Start(exp m.KMapImpl, context *SwagContext) error
}

func SwagTasksQueueNew() SwagTasksQueueImpl {

  tasks := &SwagTasksQueue{}
  tasks.Init()

  return tasks
}

func (t *SwagTasksQueue) Init() {

  t.tasks = MakeSwagTasks()
}

func (t *SwagTasksQueue) AddTask(task *SwagTask) {

  t.tasks.Push(task)
}

func (t *SwagTasksQueue) Start(exp m.KMapImpl, context *SwagContext) error {

  // why use iteration, bcs for limited searching, not for all

  var i uint

  // minify field of searching
  iter := exp.Tree().Iterable()

  // for _, enum := range exp.Tree().Enums() {

  for next := iter.Next(); next.HasNext(); next = next.Next() {

    enum := next.Enum()
    k, v := enum.Tuple()

    for i = 0; i < t.tasks.Len(); i++ {

      task, err := t.tasks.Get(i)

      if err != nil {

        fmt.Println("task failed to execute ...")
        continue
      }

      if k == task.Name {

        context.Solve(v)

        e := task.Handler(context)

        if e != nil {

          return errors.New("process failed task")
        }

        // catch var `context.Closed`
        if context.Revoke() {

          return errors.New("revoke context")
        }
      }
    }
  }

  return nil
}
