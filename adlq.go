package main

import (
	"adlq/cmd/consumer"
	"adlq/cmd/crawler"
	"adlq/cmd/producer"
	"adlq/logger"
	"flag"
	"log"
)

var mode int

const (
	Consumer = iota
	Producer
	Crawler
)

func main() {
	flag.IntVar(&mode, "mode", Crawler, "mode: 0 - consumer, 1 - producer, 2 - crawler")
	flag.Parse()

	switch mode {
	case Consumer:
		logger.Set("consumer.txt")
		consumer.Exec()
	case Producer:
		logger.Set("producer.txt")
		producer.Exec()
	case Crawler:
		logger.Set("crawler.txt")
		crawler.Exec()
	default:
		log.Fatalln("invalid mode")
	}
}
