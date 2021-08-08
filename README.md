# imdb-scraping-and-analysing

Big data storage and processing / HUST / 20201

## How to run

**Điều kiện tiên quyết**

- Docker
- Docker Compose
- MongoDB

1. Cài đặt và khởi chạy Kafka

- Khởi chạy docker-compose.yml
  ```shell
    $ cd deployment
    $ docker-compose -f docker-compose.yml up -d
  ```
- Exec vào container kafka để tạo topic
  ```shell
    $ docker exec -it kafka /bin/sh
    $ cd opt
    $ cd bitnami
    $ cd kafka
    $ cd bin
  ```
- Để tạo/xóa topic
  ```shell
  # tạo
  kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --topic movie
  # xóa
  kafka-topics.sh --delete --zookeeper zookeeper:2181 --topic movie
  ```

3. Khởi chạy Consumer

_Nhớ chỉnh sửa `kafkaUrl` và `topic` cho phù hợp với địa chỉ IP của máy và topic bạn đặt ở trên_

```shell
$ go run adlq -mode=0
```

4. Khởi chạy Producer

**File: `cmd/producer/producer.go`**

* _Nhớ chỉnh sửa `kafkaUrl` và `topic` cho phù hợp với địa chỉ IP của máy và topic bạn đặt ở trên_

* _Có thể chỉnh sửa `Parallelism` và `Delay`, tuy nhiên, hãy lịch sự nếu không muốn bị chặn_

* _Chỉ hỗ trợ crawl dữ liệu `Feature Film` tại [IMDb: Advanced Title Search](https://www.imdb.com/search/title/) nên hãy
  nhớ chỉ tick ở `Feature Film`, chọn `Search`, sau đó copy URL đó vào `c1.Visit`_

```shell
$ go run adlq -mode=1
```

5. Cài đặt và khởi chạy MongoDB Charts

- File `.yml` đã có sẵn trong `deployment`
- Làm theo hướng dẫn tại [Install MongoDB Charts](https://docs.mongodb.com/charts/current/installation/)
- Sau khi hoàn tất, mở trình duyệt và truy cập `localhost:9699`