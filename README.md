# heads-up

heads-up is a server that uses [Tile38](https://tile38.com) to store the location of the International Space Station and send notifications of when it will be over head.

## Requirements

* Tile38

## Build

```sh
make
```

## Run

```sh
heads-up -t localhost:9851 -i 5 -s 5000 -l 33.4484,112.0740 -e https://locahost:9999/webhook
```

## Notifications

Tile38 supports the following notifications, amungst others:

* http
* https
* grpc
* redis
* disque
* kafka
* amqp
* mqtt
* sqs
* nats

A URL for the above notification types can be provied with the `-e` flag.
