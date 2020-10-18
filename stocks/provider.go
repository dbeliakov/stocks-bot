package stocks

type Provider interface {
	CurrentPrice(symbol string) (float64, error)
}