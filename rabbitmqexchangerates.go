package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/aleasoluciones/simpleamqp"
)

func countRoutingKeys(routingKeys chan string) {

	ticker := time.NewTicker(60 * time.Second)
	routingKeysCounter := make(map[string]int)
	for {
		select {
		case routingKey := <-routingKeys:
			routingKeysCounter[routingKey]++
		case <-ticker.C:
			keys := make([]string, 0, len(routingKeysCounter))
			for routingKey := range routingKeysCounter {
				keys = append(keys, routingKey)
			}
			sort.Strings(keys)
			for _, key := range keys {
				log.Println(fmt.Sprintf("#/min %-4d %s", routingKeysCounter[key], key))
			}
			log.Println()
		}
	}
}

func main() {
	exchange := flag.String("exchange", "events", "exchange name")
	amqpuri := flag.String("amqpuri", "amqp://guest:guest@localhost/", "AMQP connection uri")
	flag.Parse()

	amqpConsumer := simpleamqp.NewAmqpConsumer(*amqpuri)
	messages := amqpConsumer.Receive(*exchange, []string{"#"}, "",
		simpleamqp.QueueOptions{Durable: false, Delete: true, Exclusive: true},
		30*time.Minute)

	routingKeys := make(chan string, 1024)

	go countRoutingKeys(routingKeys)

	for message := range messages {
		routingKeys <- message.RoutingKey
	}

}
