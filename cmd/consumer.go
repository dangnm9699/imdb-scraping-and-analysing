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
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync/atomic"
	"time"
)

var (
	kafkaBrokers       string
	kafkaConsumerTopic string

	mongoUri        string
	mongoDatabase   string
	mongoCollection string
	mongoUpsert     = true
	mongoTimeout    = 10 * time.Second

	mongoColl *mongo.Collection

	kafkaReader *kafka.Reader

	consumerCounter int64
)

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Start consumer",
	Long:  `Start consumer that get message from kafka and store into mongodb.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check required values
		if len(kafkaBrokers) == 0 {
			logger.Error.Fatalf("kafka broker(s) not provided\n")
		}
		if len(kafkaConsumerTopic) == 0 {
			logger.Error.Fatalf("kafka topic not provided\n")
		}
		// Connect to mongodb
		connect2MongoDb()
		// Start consumer
		startConsumer()
		// Start consuming
		ctx := context.Background()
		for {
			msg, err := kafkaReader.ReadMessage(ctx)
			if err != nil {
				logger.Debug.Printf("cannot read message: %v\n", err)
				continue
			}
			go store(string(msg.Key), string(msg.Value))
		}
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)

	consumerCmd.Flags().StringVar(&kafkaBrokers, "brokers", "", "kafka broker address list, separated by comma")
	consumerCmd.Flags().StringVar(&kafkaConsumerTopic, "topic", "", "kafka topic")

	consumerCmd.Flags().StringVar(&mongoUri, "mongo-uri", "mongodb://localhost:27017", "mongodb connection uri")
	consumerCmd.Flags().StringVar(&mongoDatabase, "mongo-db", "imdb", "specify mongodb database")
	consumerCmd.Flags().StringVar(&mongoCollection, "mongo-co", "movies", "specify mongodb collection")
}

func connect2MongoDb() {
	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		logger.Error.Fatalf("connect to mongodb: %v\n", err)
	}
	logger.Debug.Printf("connect to mongodb: ok")
	mongoColl = client.Database(mongoDatabase).Collection(mongoCollection)
}

func startConsumer() {
	brokers := strings.Split(kafkaBrokers, ",")
	kafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    kafkaConsumerTopic,
		MinBytes: 1,
		MaxBytes: 10e8,
	})
	if kafkaReader == nil {
		logger.Error.Fatalf("start consumer failed\n")
	}
}

func store(key, raw string) {
	var movie model.Movie
	var count = atomic.AddInt64(&consumerCounter, 1)
	if err := json.Unmarshal([]byte(raw), &movie); err != nil {
		logger.Debug.Printf("no.%d movie{url=%s} unmarshalled failed\n", count, key)
		return
	}
	if _, err := mongoColl.ReplaceOne(
		context.TODO(),
		bson.M{"url": movie.Url},
		movie,
		&options.ReplaceOptions{Upsert: &mongoUpsert},
	); err != nil {
		logger.Debug.Printf("no.%d movie{url=%s} saved failed\n", count, movie.Url)
	} else {
		logger.Info.Printf("no.%d movie{url=%s} saved\n", count, movie.Url)
	}
}
