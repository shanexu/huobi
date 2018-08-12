package trade

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"

	"github.com/shanexu/huobi/model"
)

var prevSellResult model.SellResult

func SellTrade() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case now := <-ticker.C:
			r, e := fetchSell()
			if e != nil {
				log.Printf("error: %s\n", e)
			} else {
				if reflect.DeepEqual(&prevSellResult, &r) {
					log.Println("no change")
				} else {
					log.Printf("ok: %+v\n", r)
					processSell(&r, now)
					prevSellResult = r
				}
			}
		}
	}
}

func fetchSell() (model.SellResult, error) {
	var result model.SellResult
	req, _ := http.NewRequest(http.MethodGet, "https://otc-api.huobi.com/v1/data/trade-market?country=37&currency=1&payMethod=0&currPage=1&pageSize=100&coinId=1&tradeType=sell&blockType=general&online=1&"+fmt.Sprintf("t=%d", time.Now().Nanosecond()), nil)
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
	defer rep.Body.Close()

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

func processSell(result *model.SellResult, now time.Time) {
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
