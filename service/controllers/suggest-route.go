package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"strconv"
)

var (
	//Chain is the global block chain used for processing
	Chain []EventBlock
	//PatientIDsMatch global patient IDs matching the first step
	PatientMatch []PatientIDsMatch
)

//PatientIDsMatch struct to match the patient id
type PatientIDsMatch struct {
	ID int
	Success bool
	Medicine string
}

//Event contains event details
type Event struct {
	PatientID int      `json:"PatientID"`
	Event     string   `json:"Event"`
	Medicine  []string `json:"Medicine"`
	TimeSFO   int      `json:"TimeSFO"`
	Success   bool     `json:"Success"`
}

//EventBlock contains event details block as used in chain
type EventBlock struct {
	EventID int
	PatientID int
	Event     string
	Medicine  []string
	TimeSFO   int
	Success   bool
	Hash	  string
}

//SuccessRate contains the success rate
type SuccessRate struct {
	Event       string  `json:"Event"`
	Medicine    string  `json:"Medicine"`
	SuccessProp float64 `json:"Probability"`
}

//SuggestHandler suggests the effective medicine for the disease and the success rate
func SuggestHandler(w http.ResponseWriter, r *http.Request) {

	// prevent CORS error
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	// patientID := r.FormValue("patientid")
	event := r.FormValue("event")
	data := Chain
	var list []EventBlock

	for i := 0; i < len(data); i++ {
		if strings.Compare(strings.ToLower(event), strings.ToLower(data[i].Event)) == 0 {
			list = append(list, data[i])
			inst := PatientIDsMatch{
				ID: data[i].PatientID,
				Success: data[i].Success,
				Medicine: data[i].Medicine[0],
			}
			PatientMatch = append(PatientMatch, inst)
		}
	}
	fmt.Println("list is *****************")
	fmt.Println(list)

	if len(list) != 0 {
		var SuccessList []SuccessRate

		for i := 0; i < len(list); i++ {
			count := 0
			for j := 0; j < len(SuccessList); j++ {
				if strings.Compare(strings.ToLower(list[i].Medicine[0]), strings.ToLower(SuccessList[j].Medicine)) == 0 {

					if list[i].Success == true {
						fmt.Println(SuccessList[j].SuccessProp)
						SuccessList[j].SuccessProp = SuccessList[j].SuccessProp + 1
					} 
					// else if list[i].Success == false {
					// 	// fmt.Println(SuccessList[j].SuccessProp)
					// 	// SuccessList[j].SuccessProp = (SuccessList[j].SuccessProp - 1)
					// }
					count++

				}
			}
			if count == 0 {
				var suclist SuccessRate
				suclist = SuccessRate{
					Event:    list[i].Event,
					Medicine: list[i].Medicine[0],
				}

				if list[i].Success == true {
					suclist.SuccessProp = 1 
				} 
				// else if list[i].Success == false {
				// 	suclist.SuccessProp = -1 
				// }
				SuccessList = append(SuccessList, suclist)
			}
		}

		for i := 0; i < len(SuccessList); i++ {
			j := fmt.Sprintf("%.2f", SuccessList[i].SuccessProp/float64(len(list)))
			value, _ := strconv.ParseFloat(j, 64)
			SuccessList[i].SuccessProp = value
		}

		j, err := json.Marshal(SuccessList)
		if err != nil {
			panic(err)
		}
		// Learning(PatientMatch)
		w.Write(j)

	} else {
		w.Write([]byte(`null`))
	}

}
