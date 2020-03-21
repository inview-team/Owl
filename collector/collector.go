package main

import (
	"context"
	"fmt"
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

var opcserver = "opc.tcp://localhost:30329"
var database = "tcp://127.0.0.1:9000?debug=true"

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

func run(node string) {
	var (
		endpoint = opcserver
		policy   = "None"
		mode     = "None"
		nodeID   = node
		interval = opcua.DefaultSubscriptionInterval.String()
	)
	fmt.Println("running! ", node)

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

	log.Print("*", ep.SecurityPolicyURI, ep.SecurityMode)

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
		fmt.Printf("error: sub=%d err=%s", sub.SubscriptionID(), err.Error())
	})

	go startCallbackSub(ctx, m, subInterval, 0, nodeID)

	go startChanSub(ctx, m, subInterval, 0, nodeID)

	<-ctx.Done()
}

func startCallbackSub(ctx context.Context, m *monitor.NodeMonitor, interval, lag time.Duration, nodes ...string) {
	sub, err := m.Subscribe(
		ctx,
		&opcua.SubscriptionParameters{
			Interval: interval,
		},
		func(s *monitor.Subscription, msg *monitor.DataChangeMessage) {
			if msg.Error != nil {
				log.Printf("[callback] sub=%d error=%s", s.SubscriptionID(), msg.Error)
			} else {
				fmt.Printf("[callback] sub=%d ts=%s node=%s value=%v", s.SubscriptionID(), msg.SourceTimestamp.UTC().Format(time.RFC3339), msg.NodeID, msg.Value.Value())
	/*
				tx := db.MustBegin()
				tx.MustExec("INSERT INTO metrics (name, timestamp, value) VALUES ($1, $2, $3)", msg.NodeID, msg.SourceTimestamp.UTC().Format(time.RFC3339), msg.Value.Value())
				tx.Commit()
*/
			}
			time.Sleep(lag)
		},
		nodes...)

	if err != nil {
		log.Fatal(err)
	}

	defer cleanup(sub)

	<-ctx.Done()
}

func startChanSub(ctx context.Context, m *monitor.NodeMonitor, interval, lag time.Duration, nodes ...string) {
	ch := make(chan *monitor.DataChangeMessage, 16)
	sub, err := m.ChanSubscribe(ctx, &opcua.SubscriptionParameters{Interval: interval}, ch, nodes...)

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
				fmt.Printf("[channel ] sub=%d ts=%s node=%s value=%v", sub.SubscriptionID(), msg.SourceTimestamp.UTC().Format(time.RFC3339), msg.NodeID, msg.Value.Value())
/*
				tx := db.MustBegin()
				tx.MustExec("INSERT INTO metrics (name, timestamp, value) VALUES ($1, $2, $3)", msg.NodeID, msg.SourceTimestamp.UTC().Format(time.RFC3339), msg.Value.Value())
				tx.Commit()
*/
			}
			time.Sleep(lag)
		}
	}
}

func cleanup(sub *monitor.Subscription) {
	fmt.Printf("stats: sub=%d delivered=%d dropped=%d", sub.SubscriptionID(), sub.Delivered(), sub.Dropped())
	sub.Unsubscribe()
}

func main() {
	db, err := sqlx.Open("clickhouse", database)
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(schema)

	go run(nodes[0])
	<-time.After(time.Second * 5)
//	for _, node := range nodes{
//		fmt.Println(node)
//		go run(node)
//	}
}
