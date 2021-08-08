package producer

import (
	"adlq/internal"
	"adlq/logger"
	"adlq/model"
	"context"
	"log"
	"time"

	"github.com/gocolly/colly"
	"github.com/segmentio/kafka-go"
)

var (
	producer  *kafka.Writer
	kafkaUrl  string
	topic     string
	movieChan chan model.MovieMsg
)

func init() {
	movieChan = make(chan model.MovieMsg, 1000)
	kafkaUrl = "192.168.1.9:9092"
	topic = "movie"
}

func Exec() {
	go startProducer()
	startCrawler()
}

func startProducer() {
	producer = &kafka.Writer{
		Addr:         kafka.TCP(kafkaUrl),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		MaxAttempts:  0,
		BatchSize:    0,
		BatchBytes:   0,
		BatchTimeout: 0,
		ReadTimeout:  0,
		WriteTimeout: 0,
		RequiredAcks: 0,
		Async:        false,
		Completion:   nil,
		Compression:  0,
		Logger:       nil,
		ErrorLogger:  nil,
		Transport:    nil,
	}
	internal.CheckNil(producer, "failed to create new kafka writer")
	for {
		movie := <-movieChan
		msg := kafka.Message{
			Key:   []byte(movie.Url),
			Value: []byte(movie.Raw),
		}
		if err := producer.WriteMessages(context.Background(), msg); err != nil {
			logger.Debug.Printf("Movie{url=%s} written to kafka failed\n", movie.Url)
		}
	}
}

func startCrawler() {
	//var count int64 = 0
	// start producing
	c1 := colly.NewCollector()
	c2 := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
	)

	_ = c2.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2, Delay: 1000 * time.Millisecond})

	c1.OnHTML("h3.lister-item-header > a[href]", func(element *colly.HTMLElement) {
		_ = c2.Visit("https://www.imdb.com" + element.Attr("href"))
	})

	c1.OnHTML("a.lister-page-next", func(element *colly.HTMLElement) {
		_ = c1.Visit("https://www.imdb.com" + element.Attr("href"))
	})

	c2.OnHTML(`script[type="application/ld+json"]`, func(element *colly.HTMLElement) {
		logger.Debug.Printf("Crawling %s\n", element.Request.URL.String())
		log.Printf("Crawling %s\n", element.Request.URL.String())
		// log.Println(element.Text)
		movie := model.MovieMsg{
			Url: element.Request.URL.String(),
			Raw: element.Text,
		}
		movieChan <- movie
	})

	// _ = c1.Visit("https://www.imdb.com/search/title/?title_type=feature&release_date=,1910-12-31&sort=year,asc")
	_ = c1.Visit("https://www.imdb.com/search/title/?title_type=feature&release_date=1911-01-01,1920-12-31&sort=year,asc")

	c1.Wait()
	c2.Wait()
	time.Sleep(time.Second)
}
