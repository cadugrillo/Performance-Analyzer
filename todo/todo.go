package todo

import (
	"errors"
	"performance-analyzer/dbdriver"
	"sync"

	"github.com/rs/xid"
)

var (
	list        []dbdriver.Todo
	mtx         sync.RWMutex
	once1       sync.Once
	once2       bool
	extDBexists bool
)

func init() {
	once1.Do(initialiseList)
}

func initialiseList() {
	list = []dbdriver.Todo{}
	once2 = true

}

// Get retrieves all elements from the todo list
func Get(userId string) []dbdriver.Todo {

	var dblist = []dbdriver.Todo{}
	dblist = append(dblist, dbdriver.DatabaseGet(userId)...)
	return dblist

}

// Add will add a new todo based on a message
func Add(userId string, message string) string {
	t := newTodo(message)
	mtx.Lock()
	//list = append(list, t)
	//if extDBexists {
	dbdriver.DatabaseAdd(userId, t.ID, t.Message, t.Complete)
	//}
	mtx.Unlock()
	return t.ID
}

// Delete will remove a Todo from the Todo list
func Delete(userId string, id string) error {
	//location, err := findTodoLocation(id)
	//if err != nil {
	//	return err
	//}
	//removeElementByLocation(location)
	//if extDBexists {
	dbdriver.DatabaseDelete(userId, id)
	//}
	return nil
}

// Complete will set the complete boolean to true, marking a todo as
// completed
func Complete(userId string, id string) error {
	//location, err := findTodoLocation(id)
	//if err != nil {
	//	return err
	//}
	//setTodoCompleteByLocation(location)
	//if extDBexists {
	dbdriver.DatabaseComplete(userId, id)
	//}
	return nil
}

func newTodo(msg string) dbdriver.Todo {
	return dbdriver.Todo{
		ID:       xid.New().String(),
		Message:  msg,
		Complete: false,
	}
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
	list = append(list[:i], list[i+1:]...)
	mtx.Unlock()
}

func setTodoCompleteByLocation(location int) {
	mtx.Lock()
	list[location].Complete = true
	mtx.Unlock()
}

func isMatchingID(a string, b string) bool {
	return a == b
}
