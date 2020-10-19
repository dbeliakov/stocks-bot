package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/looplab/fsm"

	"github.com/dbeliakov/stocks-bot/stocks"
	"github.com/dbeliakov/stocks-bot/storage"
)

const (
	Init             = "init"
	PriceWaitSymbol  = "price_wait_symbol"
	AddWaitSymbol    = "add_wait_symbol"
	RemoveWaitSymbol = "remove_wait_symbol"
)

const (
	CmdPrice          = "price"
	CmdStart          = "start"
	CmdAdd            = "add"
	CmdRemove         = "remove"
	CmdMy             = "my"
	PriceEnterSymbol  = "price_enter_symbol"
	AddEnterSymbol    = "add_enter_symbol"
	RemoveEnterSymbol = "remove_enter_symbol"
)

var (
	supportedCommands = []string{CmdStart, CmdPrice, CmdAdd, CmdRemove, CmdMy}
)

const (
	HelpText = `Привет! Я помогу тебе следить за состоянием твооего портфеля акций.
Для взаимодействия выбери одну из доступных команд.`
	EnterSymbolText           = "Введите название символа"
	UnsupportedTransitionText = "Недопустимый ввод"
	InternalErrorText         = "Внутренняя ошибка сервиса"
	PriceText                 = "Цена символа `%s`: %.2f"
	SuccessText               = "Успешно"
	NoSymbolsText             = "Нет добавленных символов"
)

type Processor struct {
	m        *fsm.FSM
	chatID   int64
	messages chan<- Reply
	provider stocks.Provider
	storage  storage.Storage
}

type IncomingMessage struct {
	Command string
	Message string
}

type Reply struct {
	ChatID  int64
	Message string
}

func NewProcessor(provider stocks.Provider, storage storage.Storage, chatID int64, messages chan<- Reply) *Processor {
	p := &Processor{
		messages: messages,
		provider: provider,
		chatID:   chatID,
		storage:  storage,
	}
	p.m = fsm.NewFSM(
		Init,
		fsm.Events{
			{Name: CmdStart, Src: []string{Init}, Dst: Init},
			{Name: CmdPrice, Src: []string{Init}, Dst: PriceWaitSymbol},
			{Name: CmdAdd, Src: []string{Init}, Dst: AddWaitSymbol},
			{Name: CmdRemove, Src: []string{Init}, Dst: RemoveWaitSymbol},
			{Name: CmdMy, Src: []string{Init}, Dst: Init},
			{Name: PriceEnterSymbol, Src: []string{PriceWaitSymbol}, Dst: Init},
			{Name: AddEnterSymbol, Src: []string{AddWaitSymbol}, Dst: Init},
			{Name: RemoveEnterSymbol, Src: []string{RemoveWaitSymbol}, Dst: Init},
		},
		fsm.Callbacks{
			CmdStart:          p.replyWithText(HelpText),
			CmdPrice:          p.replyWithText(EnterSymbolText),
			CmdAdd:            p.replyWithText(EnterSymbolText),
			CmdRemove:         p.replyWithText(EnterSymbolText),
			AddEnterSymbol:    p.addSymbol,
			RemoveEnterSymbol: p.removeSymbol,
			CmdMy:             p.showMy,
			PriceEnterSymbol:  p.writePrice,

			"enter_state": p.saveState,
		},
	)
	state, _ := p.storage.GetState(chatID)
	if state != "" {
		p.m.SetState(state)
	}
	return p
}

func (p *Processor) Process(m IncomingMessage) {
	if m.Command != "" {
		for _, c := range supportedCommands {
			if c != m.Command {
				continue
			}
			if !p.m.Can(m.Command) {
				p.messages <- p.toReply(UnsupportedTransitionText)
				p.m.SetState(Init)
				return
			}
			if err := p.m.Event(m.Command); err != nil && !errors.As(err, &fsm.NoTransitionError{}) {
				// Unexpected code branching
				log.Printf("[ERROR] Failed to make transition: %+v", err)
				p.messages <- p.toReply(InternalErrorText)
				p.m.SetState(Init)
				return
			}
			return
		}

		p.messages <- p.toReply(UnsupportedTransitionText)
		p.m.SetState(Init)
		return
	}

	if p.m.Current() == Init {
		p.messages <- p.toReply(HelpText)
		return
	}

	transitions := p.m.AvailableTransitions()
	if len(transitions) != 1 {
		log.Printf("[ERROR] More than one transition for text")
		p.messages <- p.toReply(InternalErrorText)
		p.m.SetState(Init)
		return
	}

	if err := p.m.Event(transitions[0], m.Message); err != nil {
		p.messages <- p.toReply(InternalErrorText)
		p.m.SetState(Init)
		return
	}
}

func (p *Processor) replyWithText(text string) fsm.Callback {
	return func(_ *fsm.Event) {
		p.messages <- p.toReply(text)
	}
}

func (p *Processor) writePrice(event *fsm.Event) {
	if len(event.Args) != 1 {
		log.Printf("[ERROR] Incorrect count of args")
		p.messages <- p.toReply(InternalErrorText)
		return
	}
	symbol := strings.TrimSpace(event.Args[0].(string))
	price, err := p.provider.CurrentPrice(symbol)
	if err != nil {
		log.Printf("[ERROR] Failed to get price: %+v", err)
		p.messages <- p.toReply(InternalErrorText)
		return
	}
	p.messages <- p.toReply(fmt.Sprintf(PriceText, symbol, price))
}

func (p *Processor) addSymbol(event *fsm.Event) {
	if len(event.Args) != 1 {
		log.Printf("[ERROR] Incorrect count of args")
		p.messages <- p.toReply(InternalErrorText)
		return
	}
	symbol := strings.TrimSpace(event.Args[0].(string))
	if err := p.storage.AddSymbol(p.chatID, symbol); err != nil {
		log.Printf("[ERROR] Failed to save symbol")
		p.messages <- p.toReply(InternalErrorText)
		return
	}
	p.messages <- p.toReply(SuccessText)
}

func (p *Processor) removeSymbol(event *fsm.Event) {
	if len(event.Args) != 1 {
		log.Printf("[ERROR] Incorrect count of args")
		p.messages <- p.toReply(InternalErrorText)
		return
	}
	symbol := strings.TrimSpace(event.Args[0].(string))
	if err := p.storage.RemoveSymbol(p.chatID, symbol); err != nil {
		log.Printf("[ERROR] Failed to save symbol")
		p.messages <- p.toReply(InternalErrorText)
		return
	}
	p.messages <- p.toReply(SuccessText)
}

func (p *Processor) showMy(event *fsm.Event) {
	symbols, err := p.storage.Symbols(p.chatID)
	if err != nil {
		log.Printf("[ERROR] Failed to get symbols: %+v", err)
		p.messages <- p.toReply(InternalErrorText)
		return
	}
	if len(symbols) == 0 {
		p.messages <- p.toReply(NoSymbolsText)
		return
	}
	var text strings.Builder
	for _, symbol := range symbols {
		price, err := p.provider.CurrentPrice(symbol)
		if err != nil {
			log.Printf("[ERROR] Failed to get price: %+v", err)
			p.messages <- p.toReply(InternalErrorText)
			return
		}
		text.WriteString(fmt.Sprintf(PriceText+"\n", symbol, price))
	}

	p.messages <- p.toReply(text.String())
}

func (p *Processor) toReply(text string) Reply {
	return Reply{
		ChatID:  p.chatID,
		Message: text,
	}
}

func (p *Processor) saveState(event *fsm.Event) {
	if err := p.storage.SetState(p.chatID, event.Dst); err != nil {
		log.Printf("[ERROR] Failed to save state: %+v", err)
	}
}
