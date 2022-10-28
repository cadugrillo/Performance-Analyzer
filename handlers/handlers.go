package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	analyze_captured_data "performance-analyzer/modules/analyze-captured-data"
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

// //////////////ANALYZE SIGNALS DATA HANDLER///////////////
func GetAnalyzedDataHandler(c *gin.Context) {
	TsIntervalString := c.Param("TsInterval")
	TsInterval, _ := strconv.ParseInt(TsIntervalString, 10, 64)
	c.JSON(http.StatusOK, analyze_signals.AnalyzeData(TsInterval))
}

func CheckTelegramsHandler(c *gin.Context) {

	file, _ := os.Create("capMqttData.json")
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err := io.Copy(writer, c.Request.Body)
	writer.Flush()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "Telegrams successfully uploaded")
}

func AnalyzeCapMqttDataHandler(c *gin.Context) {
	TsIntervalString := c.Param("TsInterval")
	TsInterval, _ := strconv.ParseInt(TsIntervalString, 10, 64)
	c.JSON(http.StatusOK, analyze_captured_data.AnalyzeData(TsInterval))
}

func JsonBodyToCapMqttData(httpBody io.ReadCloser) (*[]analyze_captured_data.Telegram, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	var CapMqttData []analyze_captured_data.Telegram
	err = json.Unmarshal(body, &CapMqttData)
	if err != nil {
		fmt.Println(err.Error())
		return nil, http.StatusBadRequest, err
	}
	return &CapMqttData, http.StatusOK, nil
}
