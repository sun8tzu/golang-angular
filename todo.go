package todo

import (
	"errors"
	"sync"

	"github.com/rs/xid"
)

var (
	list []Todo
	mtx  sync.RWMutex
	once sync.Once
)

// Todo data structure for a task with a description of what to do.
type Todo struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	Complete bool   `json:"complete"`
}

func init() {
	once.Do(initializeList)
}

func initializeList() {
	list = []Todo{}
}

// Get retrieves all elements from the todo list.
func Get() []Todo {
	return list
}

// Add will add a new todo based on a message.
func Add(message string) string {
	t := newTodo(message)
	mtx.Lock()
	list = append(list, t)
	mtx.Unlock()
	return t.ID
}

func newTodo(msg string) Todo {
	return Todo{
		ID:       xid.New().String(),
		Message:  msg,
		Complete: false,
	}
}

// Delete will remove a Todo from the Todo list.
func Delete(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}
	removeElementByLocation(location)
	return nil
}

// Complete will set the complete boolean to true, marking a todo as
// completed.
func Complete(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}
	setTodoCompletelyByLocation(location)
	return nil
}

func findTodoLocation(id string) (int, error) {
	mtx.RLock()
	defer mtx.RUnlock()
	for i, t := range list {
		if isMatchingID(t.ID, id) {
			return i, nil
		}
	}
	return 0, errors.New("could not find todo based on id")
}

func removeElementByLocation(i int) {
	mtx.Lock()
	defer mtx.Unlock()
	list = append(list[:i], list[i+1:]...)
}

func setTodoCompletelyByLocation(location int) {
	mtx.Lock()
	mtx.Unlock()
	list[location].Complete = true
}

func isMatchingID(a, b string) bool {
	return a == b
}
