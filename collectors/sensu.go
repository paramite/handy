package sensu


import (
    "errors"
    "fmt"
    "log"
    "strings"
    "time"
    "github.com/streadway/amqp"
)


func report(level string, err error, msg string) {
    if err != nil {
        var handle func(string, ...interface{})
        switch level {
            case "error":
                handle = log.Fatalf
            default:
                handle = log.Printf
        }
        handle("[%s] %s: %s", strings.ToUpper(level), msg, err)
	  }
}


type SensuCollector struct {
    Host string

    Config map[string]string
    Subscription []string

    conn *amqp.Connection
    chnl *amqp.Channel
    pullQue amqp.Queue
    pushQue amqp.Queue
}


func (self *SensuCollector) Process() {
    connstr := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
        self.Config["user"],
        self.Config["password"],
        self.Config["host"],
        self.Config["port"],
        self.Config["vhost"],
    )
    conn, err := amqp.Dial(connstr)
    report("error", err, "Failed to connect to RabbitMQ")
    self.conn = conn
    defer self.conn.Close()

    self.chnl, err = conn.Channel()
    report("error", err, "Failed to open a channel")
    defer self.chnl.Close()

    // declare an exchange for this client
    err = self.chnl.ExchangeDeclare(
        fmt.Sprintf("client:%s", self.Host), // name
        "fanout",                           // type
        true,                               // durable
        false,                              // auto-deleted
        false,                              // internal
        false,                              // no-wait
        nil,                                // arguments
    )
    report("error", err, "Failed to declare exchange for client")

    // declare a queue for this client
    timestamp := time.Now().Unix()
    self.pullQue, err = self.chnl.QueueDeclare(
        fmt.Sprintf("%s-mrhandy-%s", self.Host, timestamp),  // name
        false,                                              // durable
        false,                                              // delete unused
        true,                                               // exclusive
        false,                                              // no-wait
        nil,                                                // arguments
    )
    report("error", err, "Failed to declare pull queue for client")

    // register consumer
    msgs, err := self.chnl.Consume(
        self.pullQue.Name,  // queue
        "mrhandy",      // consumer
        true,           // auto ack
        false,          // exclusive
        false,          // no local
        false,          // no wait
        nil,            // args
    )
    report("error", err, "Cannot register consumer")

    // bind client queue with subscriptions
    for _, sub := range self.Subscription {
        err := self.chnl.QueueBind(
            self.pullQue.Name, // queue name
            "",            // routing key
            sub,           // exchange
            false,
            nil,
        )
        report("warning", err,
            fmt.Sprintf("Failed to bind client queue to exchange %s", sub),
        )
    }

    forever := make(chan bool)

    go func() {
        for request := range msgs {
            report("info", errors.New("request: "), fmt.Sprintf("%s", request.Body))
        }
    }()

    <-forever
}
