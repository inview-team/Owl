package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/lytics/anomalyzer"
)

var currentAnom = make(map[string]float64)
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

func MinMax(array []float64) (float64, float64) {
	var max float64 = array[0]
	var min float64 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func average() float64 {
	var res float64
	for _, val := range currentAnom {
		res += val
	}
	return res / float64(len(currentAnom))
}

func anomalyDetect(node string) bool {
	val := getTimeSeries(node)
	min, max := MinMax(val)
	conf := &anomalyzer.AnomalyzerConf{
		Sensitivity: 0.01,
		UpperBound:  max + 1,
		LowerBound:  min - 1,
		ActiveSize:  1,
		NSeasons:    4,
		Methods:     []string{"diff", "fence", "highrank", "lowrank", "magnitude"},
	}

	anom, err := anomalyzer.NewAnomalyzer(conf, val)
	if err != nil {
		log.Fatal(err)
	}

	probability := anom.Eval()
	log.Printf("Metric: %s; Probability: %f\n", nodes[node], probability)

	currentAnom[nodes[node]] = probability
	if average() > 0.85 {
		log.Printf("ANOMALY! %f\n", probability)
		err = sendNotification(time.Now().String(), nodes[node])
		if err != nil {
			log.Fatal(err)
		}
		return true
	}
	return false
}

func main() {
	fmt.Printf("%v\n", getTimeSeries("Pressure"))

	for {
		anom := false
		for n, _ := range nodes {
			anom = anomalyDetect(n)
			if anom {
				break
			}
		}
		if anom {
			time.Sleep(5 * time.Minute)
		}
	}
}
