package swag

import (
  "PapayaNet/papaya/koala/collection"
  m "PapayaNet/papaya/koala/mapping"
  "fmt"
  "github.com/gofiber/fiber/v2"
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
  AddTask(name string, handler SwagRouteHandler)
  Start(exp m.KMapImpl, ctx *fiber.Ctx) (bool, error)
}

func SwagTasksQueueNew() SwagTasksQueueImpl {

  tasks := &SwagTasksQueue{}
  tasks.Init()

  return tasks
}

func (t *SwagTasksQueue) Init() {

  t.tasks = MakeSwagTasks()
}

func (t *SwagTasksQueue) AddTask(name string, handler SwagRouteHandler) {

  task := MakeSwagTask(name, handler)
  t.tasks.Push(task)
}

func (t *SwagTasksQueue) Start(exp m.KMapImpl, ctx *fiber.Ctx) (bool, error) {

  // why use iteration, bcs for limited searching, not for all

  var i uint

  iter := exp.Tree().Iterable()

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

        context := MakeSwagContextWithEvent(ctx, v)
        err := task.Handler(context)

        // catch var `context.Closed`
        return context.Closed, err
      }
    }
  }

  return false, nil
}
