package handlers

import (
	"encoding/json"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"

	"performance-analyzer/dbdriver"
	xlsxsignals "performance-analyzer/parse-signals"
	"performance-analyzer/todo"

	"github.com/gin-gonic/gin"
)

//////////////PARSE SIGNALS HANDLER/////////////////
func ParseSignalsHandler(c *gin.Context) {
	statusCode, err := FileBodyToExcel(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(http.StatusOK, xlsxsignals.ParseSignals())

}

// GetTodoListHandler returns all current todo items
func GetTodoListHandler(c *gin.Context) {
	userID := c.Param("userId")
	c.JSON(http.StatusOK, todo.Get(userID))
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

func FileBodyToExcel(httpBody io.ReadCloser) (int, error) {
	file, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return http.StatusBadRequest, err
	}

	ioutil.WriteFile("Signals.xlsx", file, fs.ModePerm)

	return http.StatusOK, nil
}
