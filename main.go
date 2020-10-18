package main

import (
	"flag"
	"log"

	"github.com/dbeliakov/stocks-bot/bot"
	"github.com/dbeliakov/stocks-bot/stocks"
	"github.com/dbeliakov/stocks-bot/stocks/cached"
	"github.com/dbeliakov/stocks-bot/stocks/finnhub"
	"github.com/dbeliakov/stocks-bot/storage"
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

func main() {
	var p stocks.Provider = finnhub.NewProvider(finnhubToken)
	p = cached.NewProvider(p, 1000)

	s, err := storage.NewStorage("./stocks.db")
	if err != nil {
		log.Fatalf("Failed to create storage: %+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("Failed to init storage: %+v", err)
	}

	b, err := bot.NewBot(tgToken, p, s)
	if err != nil {
		log.Fatalf("Failed to create bot API: %+v", err)
	}
	if err := b.Run(); err != nil {
		log.Fatalf("Failed to run bot: %+v", err)
	}
}
