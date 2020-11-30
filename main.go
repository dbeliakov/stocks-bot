package main

import (
	"flag"
	"log"
	"time"

	"github.com/dbeliakov/stocks-bot/bot"
	"github.com/dbeliakov/stocks-bot/metrics"
	"github.com/dbeliakov/stocks-bot/stocks"
	"github.com/dbeliakov/stocks-bot/stocks/cached"
	"github.com/dbeliakov/stocks-bot/stocks/finnhub"
	"github.com/dbeliakov/stocks-bot/storage/boltdb"
)

var (
	finnhubToken string
	tgToken      string
)

func init() {
	flag.StringVar(&finnhubToken, "finnhub-token", "", "Finnhub API token")
	flag.StringVar(&tgToken, "tg-token", "", "Telegram Bot API token")
	flag.Parse()
}

const (
	finnhubURL = "https://finnhub.io/api/v1/quote"
)

func runHeartbeat() {
	const timeStep = time.Second

	counter := metrics.NewCounter("heartbeat")
	go func() {
		for {
			counter.Inc()
			time.Sleep(timeStep)
		}
	}()
}

func main() {
	var p stocks.Provider = finnhub.NewProvider(finnhubURL, finnhubToken)
	p = cached.NewProvider(p, 1000, time.Minute)

	s, err := boltdb.NewStorage("./stocks.db")
	if err != nil {
		log.Fatalf("Failed to create storage: %+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("Failed to init storage: %+v", err)
	}

	go func() {
		if err := metrics.RunServer(":2112"); err != nil {
			log.Fatalf("Failed to start server: %+v", err)
		}
	}()
	runHeartbeat()

	b, err := bot.NewBot(tgToken, p, s)
	if err != nil {
		log.Fatalf("Failed to create bot API: %+v", err)
	}
	if err := b.Run(); err != nil {
		log.Fatalf("Failed to run bot: %+v", err)
	}
}
