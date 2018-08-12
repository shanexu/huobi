package trade

import (
	"log"
	"net/http"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

var httpClient *http.Client
var influxClient influx.Client
var cookies []*http.Cookie

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
