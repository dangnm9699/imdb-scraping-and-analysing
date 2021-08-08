package consumer

import (
	"adlq/config"
	"adlq/logger"
	"adlq/model"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
)

var (
	client   *mongo.Client
	consumer *kafka.Reader
	kafkaUrl string
	topic    string
	groupID  string
	upsert   = true
)

func init() {
	kafkaUrl = "192.168.1.9:9092"
	topic = "movie"
}

func Exec() {
	startMongoDb()
	startConsumer()
}

func startMongoDb() {
	client = config.Connect2MongoDB()
}

func startConsumer() {
	brokers := strings.Split(kafkaUrl, ",")
	consumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10,
		MaxBytes: 10e8,
	})
	// start consuming
	for {
		m, err := consumer.ReadMessage(context.Background())
		if err != nil {
			logger.Debug.Println("Error")
			continue
		}
		log.Println(string(m.Key))
		go process(string(m.Key), string(m.Value))
	}
}

func process(key, raw string) {
	var movie model.Movie
	if err := json.Unmarshal([]byte(raw), &movie); err != nil {
		logger.Debug.Printf("Movie{url=%s} unmarshalled failed\n", key)
		return
	}
	if _, err := client.Database("imdb").Collection("movies").ReplaceOne(
		context.TODO(),
		bson.M{"url": movie.Url},
		movie,
		&options.ReplaceOptions{Upsert: &upsert},
	); err != nil {
		logger.Debug.Printf("Movie{url=%s} saved failed\n", movie.Url)
	}
	logger.Info.Printf("Movie{url=%s} saved\n", movie.Url)
}
