package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"github.com/nav-inc/datetime"
)

func sess() string {
        user := "" // username from registration in quote 
        pass := "" // password from registration in quote 
	
	postDataM := map[string]interface{}{
		"Email":    user,
		"Password": pass,
	}
	postData, err := json.Marshal(postDataM)
	if err != nil {
		fmt.Printf("Cannot convert map to json: %v", err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.toplo.bg/api/auth/login", bytes.NewBuffer(postData))
	if err != nil {
		fmt.Printf("POST request didn't pass: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("there is no resp: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var cont map[string]any
	json.Unmarshal(body, &cont)
	sess := fmt.Sprintf("%s", cont["token"])
	return sess
}

func jOut(sess string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.toplo.bg/api/Stations/GetStationStandartView", nil)
	if err != nil {
		fmt.Printf("there is no resp on json request: %v", err)
	}
	var bearer = "Bearer " + sess
	req.Header.Add("Authorization", bearer)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("json response is empty: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("cannot read json response: %v", err)
	}
	var cont map[string]any
	json.Unmarshal(body, &cont)

	dateTime, err := datetime.Parse((cont["$values"].([]interface{})[0].(map[string]interface{})["measuredDate"]).(string), time.Local)
	if err != nil {
		fmt.Printf("cannot convert isotime to \"normal\" : %v", err)
	}
	fmt.Printf("Адрес : %s\n", cont["$values"].([]interface{})[0].(map[string]interface{})["name"])
	fmt.Printf("Дата и час на замерване от топлофикация : %s\n", dateTime)
	fmt.Printf("Температура околна среда(извън блока) : %v\n", cont["$values"].([]interface{})[0].(map[string]interface{})["outsideTemperature"])
	fmt.Printf("Температура топла вода на входа на блока : %v\n", cont["$values"].([]interface{})[0].(map[string]interface{})["heatmeterTEmitting"])
	fmt.Printf("Температура топла вода за парно : %v\n", cont["$values"].([]interface{})[0].(map[string]interface{})["heatingMeasuredTemperature"])
	fmt.Printf("Температура топла вода за ВиК : %v\n", cont["$values"].([]interface{})[0].(map[string]interface{})["domesticHotWaterMeasuredTemperature"])

	//fileContent, _ := json.MarshalIndent(cont["$values"].([]interface{})[0], "", "  ")
	//fmt.Printf("%s\n", fileContent)
}

func main() {
	aSession := sess()
	jOut(aSession)
}
