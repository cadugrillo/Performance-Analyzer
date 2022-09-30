package handlers

import (
	"encoding/json"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strconv"

	analyze_signals "performance-analyzer/modules/analyze-signals"

	"github.com/gin-gonic/gin"
)

// /////////////PARSE SIGNALS HANDLER/////////////////
func ParseSignalsHandler(c *gin.Context) {
	statusCode, err := FileBodyToExcel(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	parsedSignals, err := analyze_signals.ParseExcelSignals()
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(http.StatusOK, parsedSignals)

}

func FileBodyToExcel(httpBody io.ReadCloser) (int, error) {
	file, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return http.StatusBadRequest, err
	}

	ioutil.WriteFile("Signals.xlsx", file, fs.ModePerm)

	return http.StatusOK, nil
}

///////////////ENDPOINT RESPONSE HANDLER/////////////////

func EndpointResponseHandler(c *gin.Context) {
	Response, statusCode, err := JsonBodyToEndpointResponse(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	EndpointResponse, err := analyze_signals.CheckEndpointResponse(Response)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(http.StatusOK, EndpointResponse)
}

func JsonBodyToEndpointResponse(httpBody io.ReadCloser) (analyze_signals.EndpointResponse, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return analyze_signals.EndpointResponse{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	var EndpointResponse analyze_signals.EndpointResponse
	err = json.Unmarshal(body, &EndpointResponse)
	if err != nil {
		return analyze_signals.EndpointResponse{}, http.StatusBadRequest, err
	}
	return EndpointResponse, http.StatusOK, nil
}

// /////////////ANALYZE DATA HANDLER///////////////
func GetAnalyzedDataHandler(c *gin.Context) {
	TsIntervalString := c.Param("TsInterval")
	TsInterval, _ := strconv.ParseInt(TsIntervalString, 10, 64)
	c.JSON(http.StatusOK, analyze_signals.AnalyzeData(TsInterval))
}

/////////////////////////////////////////////////////

// func convertHTTPBodyToTodo(httpBody io.ReadCloser) (dbdriver.Todo, int, error) {
// 	body, err := ioutil.ReadAll(httpBody)
// 	if err != nil {
// 		return dbdriver.Todo{}, http.StatusInternalServerError, err
// 	}
// 	defer httpBody.Close()
// 	return convertJSONBodyToTodo(body)
// }

// func convertJSONBodyToTodo(jsonBody []byte) (dbdriver.Todo, int, error) {
// 	var todoItem dbdriver.Todo
// 	err := json.Unmarshal(jsonBody, &todoItem)
// 	if err != nil {
// 		return dbdriver.Todo{}, http.StatusBadRequest, err
// 	}
// 	return todoItem, http.StatusOK, nil
// }
