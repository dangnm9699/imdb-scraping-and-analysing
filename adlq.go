package main

import (
	"adlq/cmd/consumer"
	"adlq/cmd/producer"
	"adlq/logger"
	"flag"
	"log"
)

var mode int

const (
	Consumer = iota
	Producer
)

func main() {
	flag.IntVar(&mode, "mode", Producer, "mode: 0 - consumer, 1 - producer")
	flag.Parse()

	switch mode {
	case Consumer:
		logger.Set("consumer.txt")
		consumer.Exec()
	case Producer:
		logger.Set("producer.txt")
		producer.Exec()
	default:
		log.Fatalln("invalid mode")
	}
}
