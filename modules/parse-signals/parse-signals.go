package parsesignals

import (
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ParsedSignals struct {
	//Aggregation string   `json:"aggregation"`
	SignalIds []string `json:"signalIds"`
}

type EndpointResponse struct {
	Signals []Signal `json:"signals"`
}

type Signal struct {
	SignalId       string  `json:"signalId"`
	LegacySignalId int64   `json:"legacySignalId"`
	Name           string  `json:"name"`
	Unit           string  `json:"unit"`
	Type           string  `json:"type"`
	AggregationId  string  `json:"aggregationId"`
	Values         []Value `json:"values"`
}

type Value struct {
	Timestamp uint64 `json:"timestamp"`
	Value     any    `json:"value"`
}

type AnalyzedData struct {
	Issues []Issue `json:"issues"`
}

type Issue struct {
	SignalId string `json:"signalId"`
	Message  string `json:"message"`
}

var (
	parsedSignals    ParsedSignals
	endpointResponse EndpointResponse
)

func ParseExcelSignals() (ParsedSignals, error) {

	parsedSignals = ParsedSignals{}

	f, err := excelize.OpenFile("Signals.xlsx")
	if err != nil {
		fmt.Println(err)
		return ParsedSignals{}, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("signals")
	if err != nil {
		fmt.Println(err)
		return ParsedSignals{}, err
	}

	//parsedSignals.Aggregation = "string"

	for i := 2; i <= len(rows); i++ {
		uuid, err := f.GetCellValue("Signals", fmt.Sprintf("E%d", i))
		if err != nil {
			fmt.Println(err)
			return ParsedSignals{}, err
		}
		parsedSignals.SignalIds = append(parsedSignals.SignalIds, uuid)
	}

	return parsedSignals, nil
}

func CheckEndpointResponse(response EndpointResponse) (EndpointResponse, error) {

	if response.Signals == nil {
		fmt.Println("Something went wrong")
		return EndpointResponse{}, errors.New("Something went wrong")
	}
	endpointResponse = response
	return endpointResponse, nil
}

func AnalyzeData() AnalyzedData {

	analyzedData := AnalyzedData{}

	for i := 0; i < len(parsedSignals.SignalIds); i++ {

		fmt.Println("outer loop started")

		for j := 0; j < len(endpointResponse.Signals); j++ {

			fmt.Println("inner loop started")

			if parsedSignals.SignalIds[i] == endpointResponse.Signals[j].SignalId {
				//do something
				fmt.Println("Found Signal")
				break
			}
			if j == len(endpointResponse.Signals)-1 {
				issue := Issue{}
				issue.SignalId = parsedSignals.SignalIds[i]
				issue.Message = "signal not found"
				analyzedData.Issues = append(analyzedData.Issues, issue)
				fmt.Println("Didn't Found Signal")
			}

		}
	}
	return analyzedData
}
