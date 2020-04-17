package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Alarm struct {
	Time string `json:"time"`
	Info string `json:"info"`
}

func sendNotification(time string, node string) error {
	alarm := Alarm{
		Time: time,
		Info: node,
	}
	reqBody, err := json.Marshal(alarm)
	if err != nil {
		return err
	}
	alarmAddr := os.Getenv("RESTAPI_ADDRESS") + "alarms"
	client := http.Client{}
	req, err := http.NewRequest("POST", alarmAddr, bytes.NewReader(reqBody))
	if err != nil {
		err = fmt.Errorf("Failed to send notification: %w", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "text/plain")
	req.Header.Add("Content-Encoding", "utf-8")
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("Failed to send notification: %w", err)
		return err
	}
	defer resp.Body.Close()
	log.Printf("Notification was sent!\n")

	return nil
}
