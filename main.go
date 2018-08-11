package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

var httpClient *http.Client
var influxClient influx.Client
var cookies []*http.Cookie

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go sellTrade()
	go buyTrade()

	<-sigs
	log.Println("Bye bye!")
}

func init() {
	httpClient = &http.Client{
		Timeout: time.Millisecond * 4500,
	}
	c, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "user",
		Password: "pass",
	})
	if err != nil {
		log.Fatal(err)
	}
	influxClient = c
}
