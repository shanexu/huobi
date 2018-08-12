package quote

import (
	"log"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/markcheno/go-quote"

	. "github.com/shanexu/huobi/influxdb"
)

func Quote() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case now := <-ticker.C:
			q, e := fetchQuote(now.Add(time.Minute*-5), now)
			if e != nil {
				log.Printf("fetchQuote error: %s\n", e)
				continue
			}
			if e := processQuote(q); e != nil {
				log.Printf("processQuote error: %s\n", e)
				continue
			}
		}
	}
}

func fetchQuote(start time.Time, end time.Time) (quote.Quote, error) {
	return quote.NewQuoteFromGdax("BTC-USD", start.Format("2006-01-02 15:04"), end.Format("2006-01-02 15:04"), quote.Min1)
}

func processQuote(q quote.Quote) error {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: "huobi",
	})
	if err != nil {
		return err
	}

	for i, _ := range q.Date {
		fields := make(map[string]interface{}, 0)
		fields["open"] = q.Open[i]
		fields["high"] = q.High[i]
		fields["low"] = q.Low[i]
		fields["close"] = q.Close[i]
		fields["volume"] = q.Volume[i]
		pt, err := influx.NewPoint("BTC-USD", map[string]string{}, fields, q.Date[i])
		if err != nil {
			continue
		}
		bp.AddPoint(pt)
	}
	InfluxClient.Write(bp)
	return nil
}
