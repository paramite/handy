package main

import "github.com/paramite/handy/processors"


func main() {
    processor := sensu.SensuProcessor{
        Host: "mrhandy-test.internal" ,
        Config: map[string]string{
            "user": "sensu",
            "password": "sensu",
            "host": "192.168.66.102",
            "port": "5672",
            "vhost": "/sensu",
        },
        Subscription: []string{"all", "overcloud-nova-conductor", "overcloud-nova-compute"},
    }
    processor.Process()
}
