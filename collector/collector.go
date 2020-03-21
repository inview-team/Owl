package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
)

var opcserver = "opc.tcp://opc-svc:8080"
var database = "tcp://clickhouse-svc:9000?debug=true"

var nodes = []string{"ns=2;i=9", "ns=2;i=10", "ns=2;i=11", "ns=2;i=12", "ns=2;i=13", "ns=2;i=14", "ns=2;i=15", "ns=2;i=16"}

var schema = ` 
CREATE TABLE IF NOT EXISTS metrics (
    name String,
    timestamp DateTime,
    value Float64
) engine=Memory`

type Metric struct {
	Name      string    `db:"name"`
	Timestamp time.Time `db:"timestamp"`
	Value     float64   `db:"value"`
}

func run(nodeID string, db *sqlx.DB) {
	var (
		endpoint = opcserver
		policy   = "None"
		mode     = "None"
		interval = opcua.DefaultSubscriptionInterval.String()
	)

	debug.Enable = true
	subInterval, err := time.ParseDuration(interval)
	if err != nil {
		log.Fatal(err)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-signalCh
		println()
		cancel()
	}()

	endpoints, err := opcua.GetEndpoints(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	ep := opcua.SelectEndpoint(endpoints, policy, ua.MessageSecurityModeFromString(mode))
	if ep == nil {
		log.Fatal("Failed to find suitable endpoint")
	}

	opts := []opcua.Option{
		opcua.SecurityPolicy(policy),
		opcua.SecurityModeString(mode),
	}

	c := opcua.NewClient(ep.EndpointURL, opts...)
	if err := c.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	m, err := monitor.NewNodeMonitor(c)
	if err != nil {
		log.Fatal(err)
	}

	m.SetErrorHandler(func(_ *opcua.Client, sub *monitor.Subscription, err error) {
		log.Printf("error: sub=%d err=%s", sub.SubscriptionID(), err.Error())
	})

	go startCallbackSub(ctx, m, subInterval, 0, db, nodeID)
	go startChanSub(ctx, m, subInterval, 0, db, nodeID)

	<-ctx.Done()
}

func startCallbackSub(ctx context.Context, m *monitor.NodeMonitor, interval, lag time.Duration, db *sqlx.DB, node string) {
	sub, err := m.Subscribe(
		ctx,
		&opcua.SubscriptionParameters{
			Interval: interval,
		},
		func(s *monitor.Subscription, msg *monitor.DataChangeMessage) {
			if msg.Error != nil {
				log.Printf("[callback] sub=%d error=%s", s.SubscriptionID(), msg.Error)
			} else {
				log.Printf("[callback] sub=%d ts=%s node=%s value=%v", s.SubscriptionID(), msg.SourceTimestamp.UTC().Format(time.RFC3339), msg.NodeID, msg.Value.Value())

				tx := db.MustBegin()
				tx.MustExec("INSERT INTO metrics (name, timestamp, value) VALUES ($1, $2, $3)", node, msg.SourceTimestamp, msg.Value.Value())
				tx.Commit()
			}
			time.Sleep(lag)
		},
		node )

	if err != nil {
		log.Fatal(err)
	}

	defer cleanup(sub)

	<-ctx.Done()
}

func startChanSub(ctx context.Context, m *monitor.NodeMonitor, interval, lag time.Duration, db *sqlx.DB, node string) {
	ch := make(chan *monitor.DataChangeMessage, 16)
	sub, err := m.ChanSubscribe(ctx, &opcua.SubscriptionParameters{Interval: interval}, ch, node)

	if err != nil {
		log.Fatal(err)
	}

	defer cleanup(sub)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			if msg.Error != nil {
				log.Printf("[channel ] sub=%d error=%s", sub.SubscriptionID(), msg.Error)
			} else {
				log.Printf("[channel ] sub=%d ts=%s node=%s value=%v", sub.SubscriptionID(), msg.SourceTimestamp.UTC().Format(time.RFC3339), msg.NodeID, msg.Value.Value())

				tx := db.MustBegin()
				tx.MustExec("INSERT INTO metrics (name, timestamp, value) VALUES ($1, $2, $3)", node, msg.SourceTimestamp, msg.Value.Value())
				tx.Commit()
			}
			time.Sleep(lag)
		}
	}
}

func cleanup(sub *monitor.Subscription) {
	log.Printf("stats: sub=%d delivered=%d dropped=%d", sub.SubscriptionID(), sub.Delivered(), sub.Dropped())
	sub.Unsubscribe()
}

func main() {
	db, err := sqlx.Open("clickhouse", database)
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(schema)
	time.Sleep(2 * time.Second)

	for {
		go run(nodes[0], db)
		go run(nodes[1], db)
		go run(nodes[2], db)
		go run(nodes[3], db)
		go run(nodes[4], db)
		go run(nodes[5], db)
		go run(nodes[6], db)
		go run(nodes[7], db)
		<-time.After(time.Second)
	}
}
