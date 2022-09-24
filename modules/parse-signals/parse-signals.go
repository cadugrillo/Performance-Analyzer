package parsesignals

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ParsedSignals struct {
	Aggregation string   `json:"aggregation"`
	SignalIds   []string `json:"signalIds"`
}

func ParseSignals() ParsedSignals {

	parsedSignals := ParsedSignals{}

	f, err := excelize.OpenFile("Signals.xlsx")
	if err != nil {
		fmt.Println(err)
		return ParsedSignals{Aggregation: err.Error()}
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
		return ParsedSignals{Aggregation: err.Error()}
	}

	parsedSignals.Aggregation = "string"

	for i := 2; i <= len(rows); i++ {
		uuid, err := f.GetCellValue("Signals", fmt.Sprintf("E%d", i))
		if err != nil {
			fmt.Println(err)
			return ParsedSignals{Aggregation: err.Error()}
		}
		parsedSignals.SignalIds = append(parsedSignals.SignalIds, uuid)
	}

	return parsedSignals
}
