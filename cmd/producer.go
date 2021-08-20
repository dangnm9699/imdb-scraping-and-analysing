// Package cmd includes CLI commands
/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"adlq/logger"
	"adlq/model"
	"context"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/segmentio/kafka-go"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
)

var (
	kafkaBroker        string
	kafkaProducerTopic string
	kafkaWriter        *kafka.Writer

	movieChan chan model.MovieMsg

	producerCounter int64
	year            int
)

// producerCmd represents the producer command
var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "Start producer",
	Long:  `Start producer that scrape movie data by year from imdb, and produce to kafka`,
	Run: func(cmd *cobra.Command, args []string) {
		// Make MovieMsg channel
		movieChan = make(chan model.MovieMsg, 1000)
		// Start producer
		go startProducer()
		// Start crawler
		startCrawler()
	},
}

func init() {
	rootCmd.AddCommand(producerCmd)

	producerCmd.Flags().StringVar(&kafkaBroker, "broker", "", "kafka broker address")
	producerCmd.Flags().StringVar(&kafkaProducerTopic, "topic", "", "kafka topic")
	producerCmd.Flags().IntVar(&year, "year", time.Now().Year(), "specify year to scrape movies")
}

func startProducer() {
	kafkaWriter = &kafka.Writer{
		Addr:         kafka.TCP(kafkaBroker),
		Topic:        kafkaProducerTopic,
		Balancer:     &kafka.LeastBytes{},
		MaxAttempts:  0,
		BatchSize:    0,
		BatchBytes:   0,
		BatchTimeout: 0,
		ReadTimeout:  0,
		WriteTimeout: 0,
		RequiredAcks: 0,
		Async:        true,
		Completion:   nil,
		Compression:  0,
		Logger:       nil,
		ErrorLogger:  nil,
		Transport:    nil,
	}
	if kafkaWriter == nil {
		logger.Error.Fatalf("start producer failed\n")
	}
	ctx := context.Background()
	for {
		movie := <-movieChan
		msg := kafka.Message{
			Key:   []byte(movie.Url),
			Value: []byte(movie.Raw),
		}
		go func() {
			if err := kafkaWriter.WriteMessages(ctx, msg); err != nil {
				logger.Debug.Printf("cannot write message: %v\n", err)
			}
		}()
	}
}

func startCrawler() {
	// start crawling
	c1 := colly.NewCollector()
	c2 := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
	)

	_ = c2.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4, Delay: 1500 * time.Millisecond})

	c1.OnHTML("h3.lister-item-header > a[href]", func(element *colly.HTMLElement) {
		_ = c2.Visit("https://www.imdb.com" + element.Attr("href"))
	})

	c1.OnHTML("a.lister-page-next", func(element *colly.HTMLElement) {
		url := element.Attr("href")
		_ = c1.Visit("https://www.imdb.com" + url)
	})

	c2.OnHTML(`script[type="application/ld+json"]`, func(element *colly.HTMLElement) {
		logger.Debug.Printf("no.%d crawling %s\n", atomic.AddInt64(&producerCounter, 1), element.Request.URL.String())
		movie := model.MovieMsg{
			Url: element.Request.URL.String(),
			Raw: element.Text,
		}
		movieChan <- movie
	})

	_ = c1.Visit(fmt.Sprintf("https://www.imdb.com/search/title/?title_type=feature&release_date=%d-01-01,%d-12-31&sort=year,asc", year, year))

	c2.Wait()
	time.Sleep(10 * time.Second)
	logger.Debug.Println("======== FINISH SESSION ========")
}
