package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
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

///////////////ENDPOINT RESPONSE HANDLERS/////////////////

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

func GetAnalyzedDataHandler(c *gin.Context) {
	TsIntervalString := c.Param("TsInterval")
	TsInterval, _ := strconv.ParseInt(TsIntervalString, 10, 64)
	c.JSON(http.StatusOK, analyze_signals.AnalyzeData(TsInterval))
}

//////////////////////////////////////////////////////////////////////////////

///////////////MQTT TELEGRAMS HANDLERS/////////////////

func AnalyzeCapMqttDataHandler(c *gin.Context) {
	TsIntervalString := c.Param("TsInterval")
	TsInterval, _ := strconv.ParseInt(TsIntervalString, 10, 64)

	Telegrams, statusCode, err := JsonBodyToCapMqttData(c.Request.Body)
	if err != nil {
		//c.JSON(statusCode, err)
		c.JSON(statusCode, analyze_captured_data.AnalyzedData{Issues: []analyze_captured_data.Issue{{SignalId: "Internal Error", Messages: []string{err.Error()}}}})
		return
	}

	c.JSON(http.StatusOK, analyze_captured_data.AnalyzeData(Telegrams, TsInterval))
}

func AnalyzeCapMqttDbusDataHandler(c *gin.Context) {
	TsIntervalString := c.Param("TsInterval")
	TsInterval, _ := strconv.ParseInt(TsIntervalString, 10, 64)

	Telegrams, statusCode, err := JsonBodyToCapMqttDbusData(c.Request.Body)
	if err != nil {
		//c.JSON(statusCode, err)
		c.JSON(statusCode, analyze_captured_data.AnalyzedData{Issues: []analyze_captured_data.Issue{{SignalId: "Internal Error", Messages: []string{err.Error()}}}})
		return
	}

	c.JSON(http.StatusOK, analyze_captured_data.AnalyzeDbusData(Telegrams, TsInterval))
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
		return nil, http.StatusOK, err
	}
	return &CapMqttData, http.StatusOK, nil
}

func JsonBodyToCapMqttDbusData(httpBody io.ReadCloser) (*[]analyze_captured_data.DbusTelegram, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	var CapMqttDbusData []analyze_captured_data.DbusTelegram
	err = json.Unmarshal(body, &CapMqttDbusData)
	if err != nil {
		fmt.Println(err.Error())
		return nil, http.StatusOK, err
	}
	return &CapMqttDbusData, http.StatusOK, nil
}
