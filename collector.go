package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
)

var OPCSERVER = ""
var DATABASE = "tcp://127.0.0.1:9000?debug=true"

var schema = ` (
    timestamp DateTime,
    value Float64
)`

type Metric struct {
	Timestamp time.Time `db:"timestamp"`
	Value     float64   `db:"value"`
}

func monitor(node string) {
	var (
		endpoint = OPCSERVER
		policy   = "None"
		mode     = "None"
		nodeID   = node
		interval = opcua.DefaultSubscriptionInterval.String()
	)

	debug.Enable = true
	subInterval, err := time.ParseDuration(*interval)
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

	endpoints, err := opcua.GetEndpoints(*endpoint)
	if err != nil {
		log.Fatal(err)
	}

	ep := opcua.SelectEndpoint(endpoints, *policy, ua.MessageSecurityModeFromString(*mode))
	if ep == nil {
		log.Fatal("Failed to find suitable endpoint")
	}

	log.Print("*", ep.SecurityPolicyURI, ep.SecurityMode)

	opts := []opcua.Option{
		opcua.SecurityPolicy(*policy),
		opcua.SecurityModeString(*mode),
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

	go startCallbackSub(ctx, m, subInterval, 0, *nodeID)

	go startChanSub(ctx, m, subInterval, 0, *nodeID)

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
				log.Printf("[callback] sub=%d ts=%s node=%s value=%v", s.SubscriptionID(), msg.SourceTimestamp.UTC().Format(time.RFC3339), msg.NodeID, msg.Value.Value())
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
				log.Printf("[channel ] sub=%d ts=%s node=%s value=%v", sub.SubscriptionID(), msg.SourceTimestamp.UTC().Format(time.RFC3339), msg.NodeID, msg.Value.Value())
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
	db, err := sqlx.Open("clickhouse", DATABASE)
	if err != nil {
		log.Fatal(err)
	}
	metricTypes := []string{"Pressure", "Humidity", "RoomTemp", "WorkTemp", "FluidFlow", "Mass", "PH", "CO2"}

	/*
		for _, mt := range metricTypes {
			metricSchema := "CREATE TABLE " + mt + schema
			db.MustExec(metricSchema)
		}
	*/
}
