package main

import (
        "bytes"
        "log"
        "strings"
        "time"

        _ "github.com/ClickHouse/clickhouse-go"
        "github.com/jmoiron/sqlx"
        "github.com/streadway/amqp"
)

var database = "tcp://clickhouse-svc:9000?debug=true"
var schema = `
CREATE TABLE IF NOT EXISTS metrics (
    name String,
    timestamp String,
    value String
) engine=Memory`

func failOnError(err error, msg string) {
        if err != nil {
                log.Fatalf("%s: %s", msg, err)
        }
}

func main() {

        conn, err := amqp.Dial("amqp://guest:guest@rabbitmq-svc:5672/")
        failOnError(err, "Failed to connect to RabbitMQ")
        defer conn.Close()

        db, err := sqlx.Open("clickhouse", database)
        if err != nil {
                log.Fatal(err)
        }
        db.MustExec(schema)

        ch, err := conn.Channel()
        failOnError(err, "Failed to open a channel")
        defer ch.Close()

        q, err := ch.QueueDeclare(
                "metrics", // name
                true,      // durable
                false,     // delete when unused
                false,     // exclusive
                false,     // no-wait
                nil,       // arguments
        )
        failOnError(err, "Failed to declare a queue")

        err = ch.Qos(
                1,     // prefetch count
                0,     // prefetch size
                false, // global
        )
        failOnError(err, "Failed to set QoS")

        msgs, err := ch.Consume(
                q.Name, // queue
                "",     // consumer
                false,  // auto-ack
                false,  // exclusive
                false,  // no-local
                false,  // no-wait
                nil,    // args
        )
        failOnError(err, "Failed to register a consumer")

        forever := make(chan bool)

        go func() {
                for d := range msgs {
                        log.Printf("Received a message: %s", d.Body)

                        // Send metric to Clickhouse
                        msg := strings.Split(string(d.Body), "; ")
                        tx := db.MustBegin()
                        tx.MustExec("INSERT INTO metrics (name, timestamp, value) VALUES ($1, $2, $3)", msg[0], msg[1], msg[2])
                        err := tx.Commit()
                        if err != nil {
                                failOnError(err, "Failed to send metrics to Clickhouse")
                        }

                        dot_count := bytes.Count(d.Body, []byte("."))
                        t := time.Duration(dot_count)
                        time.Sleep(t * time.Second)
                        log.Printf("Done")

                        err = d.Ack(false)
                        if err != nil {
                                failOnError(err, "Failed to ack message")
                        }
                }
        }()

        log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
        <-forever
}
