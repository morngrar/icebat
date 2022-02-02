package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type responseType struct {
	StatusGetRSP struct {
		WebHeader struct {
			BatteryChargeLevelPercentage int
			BatteryChargingState         bool
			BatteryPresent               bool
		}
	}
}

var chargingMap = map[bool]string{
	true:  "charging",
	false: "discharging",
}

var presenceMap = map[bool]string{
	true:  "battery",
	false: "AC",
}

func main() {
	url := "http://192.168.0.1/fcgi/cusconf.fcgi"
	jsonString := []byte(`{"StatusGetReq":{"Oper":"WebHeader"}}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	if err != nil {
		log.Fatalf("Unable to create request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Got error performing request: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}

	res := responseType{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatalf("Error unmarshaling response: %s", err)
	}

	s := fmt.Sprintf(
		"%d\t%s\t%s",
		res.StatusGetRSP.WebHeader.BatteryChargeLevelPercentage,
		chargingMap[res.StatusGetRSP.WebHeader.BatteryChargingState],
		presenceMap[res.StatusGetRSP.WebHeader.BatteryPresent],
	)

	fmt.Println(s)
}
