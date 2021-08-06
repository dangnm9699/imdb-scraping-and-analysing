package config

import (
	"adlq/internal"
	"adlq/logger"
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var (
	client *mongo.Client
)

type MongoConfig struct {
	Host     string `env:"MONGO_HOST" envDefault:"localhost"`
	Port     string `env:"MONGO_PORT" envDefault:"27017"`
	Database string `env:"MONGO_DB" envDefault:"imdb"`
}

func (conf *MongoConfig) getURI() string {
	return fmt.Sprintf("mongodb://%s:%s/%s", conf.Host, conf.Port, conf.Database)
}

func getMongoConfig() *MongoConfig {
	config := &MongoConfig{}
	internal.CheckErr(env.Parse(config), "while mongodb config")
	return config
}

func GetDB() *mongo.Client {
	return client
}

func Connect2MongoDB() {
	config := getMongoConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.getURI()))
	internal.CheckErr(err, "while connect to mongodb")
	err = client.Ping(ctx, readpref.Primary())
	internal.CheckErr(err,"while checking connection to mongodb")
	logger.Debug.Println("MongoDB connected")
}
