package swag

import (
	"PapayaNet/papaya/koala/collection"
	m "PapayaNet/papaya/koala/mapping"
	"fmt"
	"log"

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
	AddTask(task *SwagTask)
	Start(exp m.KMapImpl, ctx *fiber.Ctx) bool
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

func (t *SwagTasksQueue) Start(exp m.KMapImpl, ctx *fiber.Ctx) bool {

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

				context := MakeSwagContextWithEvent(ctx, v)
				e := task.Handler(context)

				if e != nil {

					log.Fatal(e)

					// stopped process response
					return true
				}

				// catch var `context.Closed`
				if context.Override() {

					return true
				}
			}
		}
	}

	return false
}
