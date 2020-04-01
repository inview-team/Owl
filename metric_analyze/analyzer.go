package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Metric struct {
	Node      string `json:"node"`
	Timestamp string `json:"timestamp"`
	Value     string `json:"value"`
}

type Settings struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	From int    `json:"from"`
	To   int    `json:"to"`
}

type Alarm struct {
	Time string `json:"time"`
	Info string `json:"info"`
}

var settingsAddr = "http://127.0.0.1:1337/settings?metric="
var alarmAddr = "http://127.0.0.1:1337/alarms"

var nodes = map[string]string{
	"ns=2;i=9":  "Pressure",
	"ns=2;i=10": "Humidity",
	"ns=2;i=11": "Room Temperature",
	"ns=2;i=12": "Working area Temperature",
	"ns=2;i=13": "pH",
	"ns=2;i=14": "Weight",
	"ns=2;i=15": "Fluid flow",
	"ns=2;i=16": "CO2",
}

func (mt *Metric) checkForAnomalies() {
	// Get current node
	currNode, ok := nodes[mt.Node]
	if !ok {
		err := fmt.Errorf("Invalid node value: %s\n", mt.Node)
		log.Fatal(err)
	}
	addr := settingsAddr + currNode

	// Get settings for current node
	settings, err := http.Get(addr)
	if err != nil {
		err = fmt.Errorf("Failed to get settings: %w", err)
		log.Fatal(err)
	}
	settingsBytes, err := ioutil.ReadAll(settings.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse settings and get limit values for metrics
	var limits = Settings{}
	err = json.Unmarshal(settingsBytes, &limits)
	if err != nil {
		err = fmt.Errorf("Failed to parse body: %s\n %w", settings.Body, err)
		log.Fatal(err)
	}

	// Get current metric value
	currValue, err := strconv.ParseFloat(mt.Value, 64)
	if err != nil {
		err = fmt.Errorf("Failed to parse value: %s\n %w", mt.Value, err)
		log.Fatal(err)
	}
	from := float64(limits.From)
	to := float64(limits.To)

	// Compare with metric value
	if (currValue > to) || (currValue < from) {
		// Send notification
		alarm := Alarm{
			Time: mt.Timestamp,
			Info: mt.Node,
		}

		reqBody, err := json.Marshal(alarm)
		if err != nil {
			log.Fatal(err)
		}

		client := http.Client{}
		req, err := http.NewRequest("POST", alarmAddr, bytes.NewReader(reqBody))
		if err != nil {
			err = fmt.Errorf("Failed to send notification: %w", err)
			log.Fatal(err)
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "text/plain")
		req.Header.Add("Content-Encoding", "utf-8")

		resp, err := client.Do(req)
		if err != nil {
			err = fmt.Errorf("Failed to send notification: %w", err)
			log.Fatal(err)
		}

		defer resp.Body.Close()

		log.Printf("Notification was sent! Response: %v\n", resp)
	}
}

func getMetrics(c *gin.Context) {
	body := c.Request.Body
	var metric = Metric{}
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bodyBytes, &metric)
	if err != nil {
		err = fmt.Errorf("Failed to parse body: %s\n %w", body, err)
		log.Fatal(err)
	}
	metric.checkForAnomalies()

	c.JSON(200, 0)
}

func main() {
	//	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/metrics", getMetrics)
	err := r.Run(":31337")
	if err != nil {
		log.Fatal(err)
	}
}
