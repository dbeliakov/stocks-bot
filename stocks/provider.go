package stocks

//go:generate mockgen -destination mock/provider.go -package mock . Provider
type Provider interface {
	CurrentPrice(symbol string) (float64, error)
}
