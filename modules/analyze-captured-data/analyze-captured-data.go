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
	telegrams        []Telegram
	analyzedData     AnalyzedData
	analysisStatus   string = "Not Running" // Not Running, Running, Finished, Aborted
	analysisRunning  bool   = false
	abortAnalysis    bool   = false
)

func CheckTelegrams(tlgs []Telegram) (string, error) {

	telegrams = tlgs
	return "Telegrams successfully uploaded", nil
}

func StartAnalysis(TsInterval int64) string {

	if analysisRunning {
		return GetAnalysisStatus()
	}

	analysisStatus = "Running"
	analysisRunning = true
	abortAnalysis = false

	go AnalyzeData(TsInterval)

	return GetAnalysisStatus()
}

func GetAnalysisStatus() string {
	return analysisStatus
}

func GetAnalysisResult() AnalyzedData {
	return analyzedData
}

func AbortAnalysis() string {
	abortAnalysis = true
	time.Sleep(2 * time.Second)
	return "Abort Analysis requested!"
}

func AnalyzeData(TsInterval int64) {

	analyzedData = AnalyzedData{}
	errorFlag := false
	issue := Issue{}
	checkedIds := []string{}
	idAlreadyChecked = false

	for i := 0; i < len(telegrams); i++ {

		if abortAnalysis == true {
			break
		}

		if strings.Contains(telegrams[i].Topic, "timeseries_json_generic") {

			fmt.Println("Starting analysis with telegram Seq = ", telegrams[i].Payload.Payload.Seq)

			for j := 0; j < len(telegrams[i].Payload.Payload.Vals); j++ {

				for p := 0; p < len(checkedIds); p++ {
					if telegrams[i].Payload.Payload.Vals[j].Id == checkedIds[p] {
						idAlreadyChecked = true
					}
				}

				if idAlreadyChecked == false {

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
								checkedIds = append(checkedIds, id)
							}
						}
					}
					if errorFlag {

						analyzedData.Issues = append(analyzedData.Issues, issue)
						issue = Issue{}
						errorFlag = false
					}
				}
				idAlreadyChecked = false
			}
		}
	}
	analysisRunning = false
	if abortAnalysis {
		analysisStatus = "Aborted"
	} else {
		analysisStatus = "Finished"
	}
	abortAnalysis = false
}
