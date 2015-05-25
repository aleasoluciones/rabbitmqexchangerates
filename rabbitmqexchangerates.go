package main

import (
	"flag"
	"fmt"
	"log"
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
			for routingKey, eventsPerMinute := range routingKeysCounter {
				log.Println(fmt.Sprintf("#/min %-4d %s", eventsPerMinute, routingKey))
				routingKeysCounter[routingKey] = 0
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
