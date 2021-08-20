# imdb-scraping-and-analysing

Big data storage and processing / HUST / 20201

## How to run

**Prerequisites**

- Docker
- Docker Compose
- MongoDB

1. Kafka Setup

- Run docker-compose.yml
  ```shell
    $ cd deployment
    $ docker-compose -f docker-compose.yml up -d
  ```
- Exec into kafka container to create topic
  ```shell
    $ docker exec -it kafka /bin/sh
    $ cd opt
    $ cd bitnami
    $ cd kafka
    $ cd bin
  ```
- Create/Delete topic
  ```shell
  # tạo
  kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --topic movie
  # xóa
  kafka-topics.sh --delete --zookeeper zookeeper:2181 --topic movie
  ```

3. Consumer Setup

* Usage:

    ```text
    Start consumer that get message from kafka and store into mongodb.
    
    Usage:
      adlq consumer [flags]
    
    Flags:
          --brokers string     kafka broker address list, separated by comma
      -h, --help               help for consumer
          --mongo-co string    specify mongodb collection (default "movies")
          --mongo-db string    specify mongodb database (default "imdb")
          --mongo-uri string   mongodb connection uri (default "mongodb://localhost:27017")
          --topic string       kafka topic
    
    Global Flags:
          --config string   config file (default is $HOME/.adlq.yaml)
    ```

* Run:

    ```shell
    go run main.go consumer [flags]
    ```

4. Producer Setup

* Usage:

    ```text
    Start producer that scrape movie data by year from imdb, and produce to kafka
    
    Usage:
      adlq producer [flags]
    
    Flags:
          --broker string   kafka broker address
      -h, --help            help for producer
          --topic string    kafka topic
          --year int        specify year to scrape movies (default 2021)
    
    Global Flags:
          --config string   config file (default is $HOME/.adlq.yaml)
    ```

* Run

    ```shell
    go run main.go producer [flags]
    ```

5. MongoDB Charts Setup

- File `.yml` available in `deployment`
- Follow the instruction at [Install MongoDB Charts](https://docs.mongodb.com/charts/current/installation/)
- Open browser at `localhost:9699`