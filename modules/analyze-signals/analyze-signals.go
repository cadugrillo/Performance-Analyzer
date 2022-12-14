package analyze_signals

import (
	"errors"
	"fmt"
	"time"

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
	Timestamp int64 `json:"timestamp"`
	Value     any   `json:"value"`
}

type AnalyzedData struct {
	Issues []Issue `json:"issues"`
}

type Issue struct {
	SignalId string   `json:"signalId"`
	Messages []string `json:"messages"`
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

func AnalyzeData(TsInterval int64) AnalyzedData {

	analyzedData := AnalyzedData{}
	errorFlag := false
	issue := Issue{}

	for i := 0; i < len(parsedSignals.SignalIds); i++ {

		issue.SignalId = parsedSignals.SignalIds[i]

		fmt.Println("Searching for Signal uuid: ", parsedSignals.SignalIds[i])

		for j := 0; j < len(endpointResponse.Signals); j++ {

			if parsedSignals.SignalIds[i] == endpointResponse.Signals[j].SignalId {

				fmt.Println("Signal uuid: ", parsedSignals.SignalIds[i], " found")

				for k := 0; k < len(endpointResponse.Signals[j].Values); k++ {

					if k == len(endpointResponse.Signals[j].Values)-1 {
						break
					}

					if endpointResponse.Signals[j].Values[k].Timestamp-endpointResponse.Signals[j].Values[k+1].Timestamp > TsInterval+500 {
						errorFlag = true
						Message := fmt.Sprintf("Missing Record between timestamp: %s and timestamp: %s", time.UnixMilli(endpointResponse.Signals[j].Values[k].Timestamp).UTC(), time.UnixMilli(endpointResponse.Signals[j].Values[k+1].Timestamp).UTC())
						issue.Messages = append(issue.Messages, Message)
						fmt.Println(Message)
					}
				}

				break
			}

			if j == len(endpointResponse.Signals)-1 {
				errorFlag = true
				Message := "Signal not Found!"
				issue.Messages = append(issue.Messages, Message)
				fmt.Println("Signal uuid: ", parsedSignals.SignalIds[i], " not found")
			}

		}
		if errorFlag {
			analyzedData.Issues = append(analyzedData.Issues, issue)
			issue = Issue{}
			errorFlag = false
		}
	}
	return analyzedData
}
