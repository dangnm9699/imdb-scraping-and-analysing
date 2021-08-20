package producer

import (
	"adlq/internal"
	"adlq/logger"
	"adlq/model"
	"context"
	"log"
	"sync/atomic"
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
	kafkaUrl = "192.168.1.5:9092"
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
		Async:        true,
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
		go func() {
			if err := producer.WriteMessages(context.Background(), msg); err != nil {
				logger.Debug.Printf("Movie{url=%s} written to kafka failed\n", movie.Url)
			} else {
				logger.Debug.Printf("Movie{url=%s} written to kafka\n", movie.Url)
			}
		}()
	}
}

func startCrawler() {
	var count int64 = 0
	// start producing
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
		log.Printf("|| %6d || Crawling %s\n", atomic.AddInt64(&count, 1), element.Request.URL.String())
		movie := model.MovieMsg{
			Url: element.Request.URL.String(),
			Raw: element.Text,
		}
		movieChan <- movie
	})

	// _ = c1.Visit("https://www.imdb.com/search/title/?title_type=feature&release_date=,1910-12-31&sort=year,asc") // ~ - 1910
	_ = c1.Visit("https://www.imdb.com/search/title/?title_type=feature&release_date=1932-01-01,1932-12-31&sort=year,asc") // 1911 - 1920

	c2.Wait()
	time.Sleep(10 * time.Second)
	log.Println("OUT")
	logger.Debug.Println("======== FINISH SESSION ========")
}
