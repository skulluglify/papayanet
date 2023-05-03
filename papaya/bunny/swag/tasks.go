package swag

import (
  "errors"
  "fmt"
  "reflect"
  "skfw/papaya/koala/collection"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
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

  // for _, enum := range exp.Tree().Enums() {

  // Queue
  // order by tasks

  for i = 0; i < t.tasks.Len(); i++ {

    task, err := t.tasks.Get(i)

    for _, enum := range exp.Tree().Enums() {

      k, v := enum.Tuple()

      if err != nil {

        fmt.Println("task failed to execute ...")
        continue
      }

      if k != task.Name {

        continue
      }

      ///////////////////////////////////

      // maybe set key but want a run task

      var played bool
      played = true

      val := pp.KIndirectValueOf(v)

      if val.IsValid() {

        ty := val.Type()

        switch ty.Kind() {

        case reflect.Bool:

          // maybe false
          played = val.Bool()
          break

        default:

          context.Solve(v)
        }
      }

      ///////////////////////////////////

      if played {

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
