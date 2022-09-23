package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"performance-analyzer/dbdriver"
	"performance-analyzer/todo"

	"github.com/gin-gonic/gin"
)

// GetTodoListHandler returns all current todo items
func GetTodoListHandler(c *gin.Context) {
	userID := c.Param("userId")
	c.JSON(http.StatusOK, todo.Get(userID))
}

// AddTodoHandler adds a new todo to the todo list
func AddTodoHandler(c *gin.Context) {
	userID := c.Param("userId")
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(statusCode, gin.H{"id": todo.Add(userID, todoItem.Message)})
}

// DeleteTodoHandler will delete a specified todo based on user http input
func DeleteTodoHandler(c *gin.Context) {
	userID := c.Param("userId")
	todoID := c.Param("id")
	if err := todo.Delete(userID, todoID); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

// CompleteTodoHandler will complete a specified todo based on user http input
func CompleteTodoHandler(c *gin.Context) {
	userID := c.Param("userId")
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	if todo.Complete(userID, todoItem.ID) != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

func convertHTTPBodyToTodo(httpBody io.ReadCloser) (dbdriver.Todo, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return dbdriver.Todo{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToTodo(body)
}

func convertJSONBodyToTodo(jsonBody []byte) (dbdriver.Todo, int, error) {
	var todoItem dbdriver.Todo
	err := json.Unmarshal(jsonBody, &todoItem)
	if err != nil {
		return dbdriver.Todo{}, http.StatusBadRequest, err
	}
	return todoItem, http.StatusOK, nil
}
