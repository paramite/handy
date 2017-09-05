package main

import "github.com/paramite/handy/collectors"


func main() {
    collector := sensu.SensuCollector{
        Host: "mrhandy-test.internal" ,
        Config: map[string]string{
            "user": "sensu",
            "password": "sensu",
            "host": "192.168.66.102",
            "port": "5672",
            "vhost": "/sensu",
        },
        Subscription: []string{"overcloud-nova-conductor", "overcloud-nova-compute"},
    }

    collector.Initialize()
    collector.Process()
    collector.Destroy()
}
