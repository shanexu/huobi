package main

import (
	"encoding/json"
	"github.com/shanexu/huobi/model"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var client *http.Client
var cookies []*http.Cookie

func main() {
	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ticker.C:
			r, e := fetch()
			if e != nil {
				log.Printf("error: %s\n", e)
			} else {
				log.Printf("ok: %s\n", r)
			}
		}
	}
}

func init() {
	client = &http.Client{
		Timeout: time.Second * 1,
	}
}

//curl 'https://api-otc.huobi.pro/v1/otc/trade/list/public?coinId=1&tradeType=1&currentPage=1&payWay=&country=&merchant=1&online=1&range=0&pageSize=100' -H 'Origin: https://otc.huobipro.com' -H 'Accept-Language: zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36' -H 'otc-language: zh-CN' -H 'Accept: application/json, text/plain, */*' -H 'Referer: https://otc.huobipro.com/' -H 'X-Requested-With: XMLHttpRequest' -H 'Connection: keep-alive'
func fetch() (model.Result, error) {
	var result model.Result
	req, _ := http.NewRequest(http.MethodGet, "https://api-otc.huobi.pro/v1/otc/trade/list/public?coinId=1&tradeType=1&currentPage=1&payWay=&country=&merchant=1&online=1&range=0&pageSize=100", nil)
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
	rep, err := client.Do(req)
	if err != nil {
		return result, err
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
