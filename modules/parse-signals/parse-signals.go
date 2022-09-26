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

var (
	parsedSignals ParsedSignals
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

func CheckEndpointResponse(endpointResponse EndpointResponse) (EndpointResponse, error) {

	if endpointResponse.Signals == nil {
		fmt.Println("Something went wrong")
		return EndpointResponse{}, errors.New("Something went wrong")
	}

	return endpointResponse, nil
}
