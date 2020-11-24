package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/dbeliakov/stocks-bot/stocks"
	"github.com/dbeliakov/stocks-bot/storage"
)

type Bot struct {
	api                  *tgbotapi.BotAPI
	stocksProvider       stocks.Provider
	states               map[int64]*Processor
	replies              chan Reply
	storage              storage.Storage
	totalMessagesCounter prometheus.Counter
}

func NewBot(apiKey string, provider stocks.Provider, storage storage.Storage) (*Bot, error) {
	impl, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create new bot: %w", err)
	}
	b := &Bot{
		api:            impl,
		stocksProvider: provider,
		states:         make(map[int64]*Processor),
		replies:        make(chan Reply),
		storage:        storage,
		totalMessagesCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "stocks_bot",
			Name:      "total_messages",
		}),
	}
	prometheus.MustRegister(b.totalMessagesCounter)
	return b, nil
}

func (b *Bot) Run() error {
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("failed to get updates: %w", err)
	}

	go b.processReplies()
	b.processMessages(updates)

	return nil
}

func (b *Bot) processMessages(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		b.totalMessagesCounter.Inc()

		m := update.Message
		processor, found := b.states[m.Chat.ID]
		if !found {
			processor = NewProcessor(b.stocksProvider, b.storage, update.Message.Chat.ID, b.replies)
			b.states[m.Chat.ID] = processor
		}
		if m.IsCommand() {
			processor.Process(IncomingMessage{
				Command: m.Command(),
				Message: m.CommandArguments(),
			})
		} else {
			processor.Process(IncomingMessage{Message: m.Text})
		}
	}
}

func (b *Bot) processReplies() {
	for r := range b.replies {
		reply := tgbotapi.NewMessage(r.ChatID, r.Message)
		if _, err := b.api.Send(reply); err != nil {
			log.Printf("[WARN] Failed to send message to chat `%d`: %+v", r.ChatID, err)
		}
	}
}
