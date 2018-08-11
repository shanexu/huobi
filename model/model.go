package model

type SellResult struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	TotalCount int         `json:"totalCount"`
	PageSize   int         `json:"pageSize"`
	CurrPage   int         `json:"currPage"`
	Data       []SellTrade `json:"data"`
	Success    bool        `json:"success"`
}

type SellTrade struct {
	ID                  int     `json:"id"`
	TradeNo             string  `json:"tradeNo"`
	Country             int     `json:"country"`
	CoinId              int     `json:"coinId"`
	TradeType           int     `json:"tradeType"`
	Merchant            int     `json:"merchant"`
	MerchantLevel       int     `json:"merchantLevel"`
	Currency            int     `json:"currency"`
	PayMethod           string  `json:"payMethod"`
	UserId              int     `json:"userId"`
	UserName            string  `json:"userName"`
	IsFixed             bool    `json:"isFixed"`
	MinTradeLimit       float64 `json:"minTradeLimit"`
	MaxTradeLimit       float64 `json:"maxTradeLimit"`
	FixedPrice          float64 `json:"fixedPrice"`
	CalcRate            float64 `json:"calcRate"`
	Price               float64 `json:"price"`
	GmtSort             int     `json:"gmtSort"`
	TradeCount          float64 `json:"tradeCount"`
	IsOnline            bool    `json:"isOnline"`
	TradeMonthTimes     int     `json:"tradeMonthTimes"`
	AppealMonthTimes    int     `json:"appealMonthTimes"`
	AppealMonthWinTimes int     `json:"appealMonthWinTimes"`
	TakerRealLevel      int     `json:"takerRealLevel"`
	TakerIsPhoneBind    bool    `json:"takerIsPhoneBind"`
	TakerTradeTimes     int     `json:"takerTradeTimes"`
	TakerLimit          int     `json:"takerLimit"`
}

type BuyResult struct {
	Code       int        `json:"code"`
	Message    string     `json:"message"`
	TotalCount int        `json:"totalCount"`
	PageSize   int        `json:"pageSize"`
	CurrPage   int        `json:"currPage"`
	Data       []BuyTrade `json:"data"`
	Success    bool       `json:"success"`
}

type BuyTrade struct {
	ID                int     `json:"id"`
	UID               int     `json:"uid"`
	UserName          string  `json:"userName"`
	MerchantLevel     int     `json:"merchantLevel"`
	CoinId            int     `json:"coinId"`
	Currency          int     `json:"currency"`
	TradeType         int     `json:"tradeType"`
	BlockType         int     `json:"blockType"`
	PayMethod         string  `json:"payMethod"`
	PayTerm           string  `json:"payTerm"`
	PayName           string  `json:"payName"`
	MinTradeLimit     float64 `json:"minTradeLimit"`
	MaxTradeLimit     float64 `json:"maxTradeLimit"`
	Price             float64 `json:"price"`
	TradeCount        float64 `json:"tradeCount"`
	IsOnline          bool    `json:"isOnline"`
	TradeMonthTimes   int     `json:"tradeMonthTimes"`
	OrderCompleteRate int     `json:"orderCompleteRate"`
	TakerLimit        int     `json:"takerLimit"`
	GmtSort           int     `json:"gmtSort"`
}
