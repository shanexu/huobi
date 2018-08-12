package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shanexu/huobi/quote"
	"github.com/shanexu/huobi/trade"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go trade.SellTrade()
	go trade.BuyTrade()
	go quote.Quote()

	<-sigs
	log.Println("Bye bye!")
}
