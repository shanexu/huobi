package trade

import (
	"net/http"
	"time"
)

var httpClient *http.Client
var cookies []*http.Cookie

func init() {
	httpClient = &http.Client{
		Timeout: time.Millisecond * 4500,
	}
}
