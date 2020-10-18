package cached

import (
	"fmt"
	"time"

	"github.com/karlseguin/ccache"

	"github.com/dbeliakov/stocks-bot/stocks"
)

type Provider struct {
	underlying stocks.Provider
	cache *ccache.Cache
}

var _ stocks.Provider = &Provider{}

func NewProvider(underlying stocks.Provider, maxSize int64) *Provider {
	return &Provider{
		underlying: underlying,
		cache:      ccache.New(ccache.Configure().MaxSize(maxSize)),
	}
}

func (p *Provider) CurrentPrice(symbol string) (float64, error) {
	item := p.cache.Get(symbol)
	if item != nil && item.TTL().Seconds() > 0 {
		return item.Value().(float64), nil
	}

	price, err := p.underlying.CurrentPrice(symbol)
	if err != nil {
		return 0, fmt.Errorf("failed to update value: %w", err)
	}
	p.cache.Set(symbol, price, /*TODO*/ time.Minute)
	return price, nil
}