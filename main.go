package main

import (
	"encoding/json"
	"fmt"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/shanexu/huobi/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

var httpClient *http.Client
var influxClient influx.Client
var cookies []*http.Cookie
var prevResult model.Result

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go sellTrade()

	<-sigs
	log.Println("Bye bye!")
}

func sellTrade() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case now := <-ticker.C:
			r, e := fetch()
			if e != nil {
				log.Printf("error: %s\n", e)
			} else {
				if reflect.DeepEqual(&prevResult, &r) {
					log.Println("no change")
				} else {
					log.Printf("ok: %+v\n", r)
					process(&r, now)
					prevResult = r
				}
			}
		}
	}
}

func buyTrade() {
	// curl 'https://api-otc.hb-otc.net/v1/otc/trade/list/public?coinId=1&tradeType=1&currentPage=1&payWay=&country=&merchant=1&online=1&range=0&pageSize=100' -H 'Origin: https://otc.huobipro.com' -H 'Accept-Language: zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36' -H 'otc-language: zh-CN' -H 'Accept: application/json, text/plain, */*' -H 'Referer: https://otc.huobipro.com/' -H 'X-Requested-With: XMLHttpRequest' -H 'Connection: keep-alive'
	// curl 'https://api-otc.hb-otc.net/v1/otc/trade/list/public?coinId=1&tradeType=0&currentPage=1&payWay=&country=&merchant=1&online=1&range=0' -H 'Cookie: acw_tc=AQAAACxVEhLdGgYAVYi8J4o98Ic2W6h/' -H 'fingerprint: 3eff0ec15b337b4930d3d390a602908a' -H 'Origin: https://otc.huobipro.com' -H 'Accept-Encoding: gzip, deflate, br' -H 'Accept-Language: zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36' -H 'otc-language: zh-CN' -H 'Accept: application/json, text/plain, */*' -H 'Referer: https://otc.huobipro.com/' -H 'X-Requested-With: XMLHttpRequest' -H 'Connection: keep-alive' -H 'token: TICKET_f0cb8e10e46d81eb250f526d6c98fd8f372acbf9e4b74a11ae77bef041b5584d' --compressed
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

func process(result *model.Result, now time.Time) {
	l := len(result.Data)
	if l > 10 {
		l = 10
	}
	sample := result.Data[0:l]
	//TOP 10 average
	amount := 0.0
	volume := 0.0
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: "huobi",
	})
	for i, _ := range sample {
		d := &sample[i]
		amount = amount + d.TradeCount*d.Price
		volume = volume + d.TradeCount
		fields := make(map[string]interface{})
		fields[fmt.Sprintf("amount%d", i+1)] = d.TradeCount * d.Price
		fields[fmt.Sprintf("volume%d", i+1)] = d.TradeCount
		fields[fmt.Sprintf("price%d", i+1)] = d.Price
		pt, err := influx.NewPoint("trade", map[string]string{}, fields, now)
		if err == nil {
			bp.AddPoint(pt)
		}
	}
	pt, err := influx.NewPoint("top10trade", map[string]string{}, map[string]interface{}{"amount": amount, "volume": volume}, now)
	if err == nil {
		bp.AddPoint(pt)
	}
	if len(sample) > 0 {
		d := &sample[0]
		pt, err := influx.NewPoint("top1trade", map[string]string{}, map[string]interface{}{"price": d.Price}, now)
		if err == nil {
			bp.AddPoint(pt)
		}
	}
	influxClient.Write(bp)
}

//curl 'https://api-otc.hb-otc.net/v1/otc/trade/list/public?coinId=1&tradeType=1&currentPage=1&payWay=&country=&merchant=1&online=1&range=0&pageSize=100' -H 'Origin: https://otc.huobipro.com' -H 'Accept-Language: zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36' -H 'otc-language: zh-CN' -H 'Accept: application/json, text/plain, */*' -H 'Referer: https://otc.huobipro.com/' -H 'X-Requested-With: XMLHttpRequest' -H 'Connection: keep-alive'
func fetch() (model.Result, error) {
	var result model.Result
	req, _ := http.NewRequest(http.MethodGet, "https://api-otc.hb-otc.net/v1/otc/trade/list/public?coinId=1&tradeType=1&currentPage=1&payWay=&country=&merchant=1&online=1&range=0&pageSize=100&"+fmt.Sprintf("t=%d", time.Now().Nanosecond()), nil)
	req.Header.Add("Origin", "https://otc.huobipro.com")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36")
	req.Header.Add("otc-language", "zh-CN")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Referer", "https://otc.huobipro.com/")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Connection", "keep-alive")
	if len(cookies) > 0 {
		for _, c := range cookies {
			req.AddCookie(c)
		}
	}
	rep, err := httpClient.Do(req)
	if err != nil {
		return result, err
	}
	if rep.Body != nil {
		defer rep.Body.Close()
	}
	data, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	if len(cookies) < 0 && len(rep.Cookies()) > 0 {
		cookies = rep.Cookies()
	}
	return result, nil
}
