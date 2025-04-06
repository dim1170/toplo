package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var apiToplo string = "https://api.toplo.bg/api"

func sess() *http.Request {

	postDataM := map[string]any{
		"Email":    os.Getenv("USER_TOPLO"), // from ENV (tested with linux env) : username from toplo registration
		"Password": os.Getenv("PASS_TOPLO"), // from ENV (tested with linux env) : password from toplo registration
	}

	postData, err := json.Marshal(postDataM)
	if err != nil {
		log.Fatal("cannot create json ", err)
	}

	req, err := http.NewRequest("POST", apiToplo+"/auth/login", bytes.NewBuffer(postData))
	if err != nil {
		log.Fatal("cannot create post request:", err)
	}
	return req
}

func jOut(sess string) (req *http.Request) {

	req, err := http.NewRequest("GET", apiToplo+"/Stations/GetStationStandartView", nil)
	if err != nil {
		log.Fatal("cannot create get request: ", err)
	}
	var bearer = "Bearer " + sess
	req.Header.Add("Authorization", bearer)
	return req
}

func httpz(valS string, token string) (body []byte) {
	client := &http.Client{}

	var req *http.Request
	if valS == "token" {
		req = sess()
	} else if valS == "jVal" {
		req = jOut(token)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("response from server is empty: ", err)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("cannot read read body: ", err)
	}

	return body

}

func main() {

	bodyT := httpz("token", "")
	var cont map[string]any
	json.Unmarshal(bodyT, &cont)

	token := fmt.Sprintf("%s", cont["token"])

	jValues := httpz("jVal", token)

	var jValue map[string]any
	json.Unmarshal(jValues, &jValue)
	wValue := jValue["$values"].([]interface{})[0].(map[string]interface{})

	t := time.Now()
	fmt.Printf("Адрес : %s\n", wValue["name"])
	fmt.Printf("Дата и час на замерване от топлофикация : %02d:%02d:%02d %02d-%02d-%d\n", t.Hour(), t.Minute(), t.Second(), t.Day(), t.Month(), t.Year())
	fmt.Printf("Температура околна среда(извън блока) : %v\n", wValue["outsideTemperature"])
	fmt.Printf("Температура топла вода на входа на блока : %v\n", wValue["heatmeterTEmitting"])
	fmt.Printf("Температура топла вода за парно : %v\n", wValue["heatingMeasuredTemperature"])
	fmt.Printf("Температура топла вода за ВиК : %v\n", wValue["domesticHotWaterMeasuredTemperature"])

}
