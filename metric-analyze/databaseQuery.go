package main

import (
	"fmt"
	"log"
	"strconv"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
)

type metric struct {
	Name      string `db:"name"`
	Timestamp string `db:"timestamp"`
	Value     string `db:"value"`
}

func getTimeSeries(targetMetric string) []float64 {
	connect, err := sqlx.Open("clickhouse", "tcp://167.172.137.177:30001")
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to connect to db: %w\n", err))
	}

	var mt []metric
	if err := connect.Select(&mt, "SELECT name, timestamp, value FROM metrics"); err != nil {
		log.Fatal(fmt.Errorf("Failed to get metric values: %w\n", err))
	}

	var res []float64
	for _, metric := range mt {
		if metric.Name == targetMetric {
			tmp, err := strconv.ParseFloat(metric.Value, 64)
			if err != nil {
				log.Fatal(fmt.Errorf("Cannot parse metric value %s\n", metric.Value))
			}
			res = append(res, tmp)
		}
	}

	return res
}
