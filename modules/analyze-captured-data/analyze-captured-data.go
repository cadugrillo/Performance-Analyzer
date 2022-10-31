package analyze_captured_data

import (
	"fmt"
	"strings"
	"time"
)

type Telegram struct {
	Topic   string       `json:"topic"`
	Payload OuterPayload `json:"payload"`
}

type OuterPayload struct {
	ClientID string       `json:"clientID"`
	Topic    string       `json:"topic"`
	Protocol string       `json:"protocol"`
	Payload  InnerPayload `json:"payload"`
}

type InnerPayload struct {
	Seq  int64    `json:"seq"`
	Vals []SigVal `json:"vals"`
}

type SigVal struct {
	Id       string `json:"id"`
	Qc       int    `json:"qc"`
	Ts       string `json:"ts"`
	Val      int64  `json:"val"`
	DataType string `json:"dataType"`
	Name     string `json:"name"`
}

type AnalyzedData struct {
	Issues []Issue `json:"issues"`
}

type Issue struct {
	SignalId string   `json:"signalId"`
	Messages []string `json:"messages"`
}

var (
	id               string
	ts               int64
	tsOffset         int64
	idAlreadyChecked bool
	analyzedData     AnalyzedData
)

func AnalyzeData(Telegrams *[]Telegram, TsInterval int64) AnalyzedData {

	telegrams := *Telegrams

	errorFlag := false
	issue := Issue{}
	checkedIds := []string{}
	idAlreadyChecked = false

	for i := 0; i < len(telegrams); i++ {

		if !strings.Contains(telegrams[i].Topic, "timeseries_json_generic") {
			continue
		}

		fmt.Println("Starting analysis with telegram Seq = ", telegrams[i].Payload.Payload.Seq)

		for j := 0; j < len(telegrams[i].Payload.Payload.Vals); j++ {

			for p := 0; p < len(checkedIds); p++ {
				if telegrams[i].Payload.Payload.Vals[j].Id == checkedIds[p] {
					idAlreadyChecked = true
					break
				}
			}

			if idAlreadyChecked {
				idAlreadyChecked = false
				continue
			}

			id = telegrams[i].Payload.Payload.Vals[j].Id
			issue.SignalId = fmt.Sprintf("%s - from Telegram Seq number = %d", telegrams[i].Payload.Payload.Vals[j].Id, telegrams[i].Payload.Payload.Seq)
			tsRFC3339, _ := time.Parse(time.RFC3339, telegrams[i].Payload.Payload.Vals[j].Ts)
			ts = tsRFC3339.UnixMilli()
			tsOffset = 200

			for k := i; k < len(telegrams); k++ {

				for l := 0; l < len(telegrams[k].Payload.Payload.Vals); l++ {

					if telegrams[k].Payload.Payload.Vals[l].Id == id {

						cmptsRFC3339, _ := time.Parse(time.RFC3339, telegrams[k].Payload.Payload.Vals[l].Ts)
						cmpts := cmptsRFC3339.UnixMilli()

						if (cmpts - ts) > TsInterval+tsOffset {

							errorFlag = true
							Message := fmt.Sprintf("Missing records near seq: %d", telegrams[k].Payload.Payload.Seq)
							issue.Messages = append(issue.Messages, Message)
							fmt.Println("Missing records from id: ", id, "near seq: ", telegrams[k].Payload.Payload.Seq)
						}
						ts = cmpts
					}
				}
			}
			checkedIds = append(checkedIds, id)
			if errorFlag {

				analyzedData.Issues = append(analyzedData.Issues, issue)
				issue = Issue{}
				errorFlag = false
			}
		}
	}

	fmt.Println("Analyzed Record IDs:", checkedIds)
	telegrams = []Telegram{}
	issue = Issue{}
	checkedIds = []string{}
	defer ClearData()

	return analyzedData
}

func ClearData() {
	analyzedData = AnalyzedData{}
}
