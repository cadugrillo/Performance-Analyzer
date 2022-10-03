package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	id       string
	ts       int64
	tsOffset int64
)

func main() {

	// analyzedData := AnalyzedData{}
	// errorFlag := false
	// issue := Issue{}

	f, err := os.Open("captured-mock-data.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var telegrams []Telegram
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&telegrams)
	if err != nil {
		panic(err)
	}

	// fmt.Println(telegrams[0].Payload.Payload.Vals[0])

	// time1, _ := time.Parse(time.RFC3339, telegrams[0].Payload.Payload.Vals[0].Ts).UnixMilli()
	// time2, _ := time.Parse(time.RFC3339, telegrams[3].Payload.Payload.Vals[0].Ts)
	// timeDiff := time2.UnixMilli() - time1.UnixMilli()

	// fmt.Println(time1)
	// fmt.Println(time2)
	// fmt.Println(timeDiff)

	for i := 0; i < len(telegrams); i++ {
		fmt.Println("Starting analysis with telegram Seq = ", telegrams[i].Payload.Payload.Seq)

		for j := 0; j < len(telegrams[i].Payload.Payload.Vals); j++ {

			id = telegrams[i].Payload.Payload.Vals[j].Id
			tsRFC3339, _ := time.Parse(time.RFC3339, telegrams[i].Payload.Payload.Vals[j].Ts)
			ts = tsRFC3339.UnixMilli()
			tsOffset = 200

			for k := i; k < len(telegrams); k++ {
				for l := 0; l < len(telegrams[k].Payload.Payload.Vals); l++ {
					if telegrams[k].Payload.Payload.Vals[l].Id == id {
						cmptsRFC3339, _ := time.Parse(time.RFC3339, telegrams[k].Payload.Payload.Vals[l].Ts)
						cmpts := cmptsRFC3339.UnixMilli()
						// fmt.Println("debug1: ", telegrams[k].Payload.Payload.Vals[l].Id) //debugging
						// fmt.Println("debug2: ", cmpts)                                   //debugging
						// fmt.Println("debug3: ", ts)                                      //debugging
						// fmt.Println("debug4: ", cmpts-ts)                                //debugging
						if (cmpts - ts) > 1000+tsOffset {
							msg := fmt.Sprintf("Possible missing records from id: %s near seq: %d", id, telegrams[k].Payload.Payload.Seq)
							fmt.Println(msg)
						}
						ts = cmpts

						// fmt.Println("debug5: ", tsOffset) //debugging

					}
				}
			}

		}
	}
}
