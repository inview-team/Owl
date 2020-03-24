package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
)

var opcserver = "opc.tcp://opc-svc:8080"
var nodes = []string{"ns=2;i=9", "ns=2;i=10", "ns=2;i=11", "ns=2;i=12", "ns=2;i=13", "ns=2;i=14", "ns=2;i=15", "ns=2;i=16"}

func failOnError(err error, msg string) {
        if err != nil {
                log.Fatalf("%s: %s", msg, err)
        }
}

func startCallbackSub(ctx context.Context, m *monitor.NodeMonitor, interval, lag time.Duration, q *amqp.Queue, ch *amqp.Channel, node string) {
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
				val := fmt.Sprintf("%v", msg.Value.Value())
				body := node + "; " + msg.SourceTimestamp.Format("2006-01-02 15:04:05") + "; " + val
				err := ch.Publish(
			                "",           // exchange
			                q.Name,       // routing key
			                false,        // mandatory
			                false,
			                amqp.Publishing{
			                        DeliveryMode: amqp.Persistent,
						ContentType:  "text/plain",
			                        Body:         []byte(body),
		                })
		        failOnError(err, "Failed to publish a message")
		        log.Printf(" [x] Sent %s", body)

			}
			time.Sleep(lag)
		},
		node )

	if err != nil {
		failOnError(err, "Failed to subscribe")
	}

	defer cleanup(sub)

	<-ctx.Done()
}

func cleanup(sub *monitor.Subscription) {
	log.Printf("stats: sub=%d delivered=%d dropped=%d", sub.SubscriptionID(), sub.Delivered(), sub.Dropped())
	err := sub.Unsubscribe()
	if err != nil {
		failOnError(err, "Failed to unsubscribe")
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq-svc:5672/")
        failOnError(err, "Failed to connect to RabbitMQ")
        defer conn.Close()
	log.Printf("Connected\n")
        ch, err := conn.Channel()
        failOnError(err, "Failed to open a channel")
        defer ch.Close()
	log.Printf("Opened channel \n")
        q, err := ch.QueueDeclare(
                "metrics", // name
                true,         // durable
                false,        // delete when unused
                false,        // exclusive
                false,        // no-wait
                nil,          // arguments
        )
        failOnError(err, "Failed to declare a queue")
	log.Printf("Declared queue %s\n", q.Name)
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

	for {
		go startCallbackSub(ctx, m, subInterval, 0, &q, ch, nodes[0])
		go startCallbackSub(ctx, m, subInterval, 0, &q, ch, nodes[1])
		go startCallbackSub(ctx, m, subInterval, 0, &q, ch, nodes[2])
		go startCallbackSub(ctx, m, subInterval, 0, &q, ch, nodes[3])
		go startCallbackSub(ctx, m, subInterval, 0, &q, ch, nodes[4])
		go startCallbackSub(ctx, m, subInterval, 0, &q, ch, nodes[5])
		go startCallbackSub(ctx, m, subInterval, 0, &q, ch, nodes[6])
		go startCallbackSub(ctx, m, subInterval, 0, &q, ch, nodes[7])
		<-ctx.Done()
		<-time.After(time.Second)
	}
}
