package storage

//go:generate mockgen -destination mock/storage.go -package mock . Storage
type Storage interface {
	Init() error
	SetState(chatID int64, state string) error
	GetState(chatID int64) (string, error)
	AddSymbol(chatID int64, symbol string) error
	RemoveSymbol(chatID int64, symbol string) error
	Symbols(chatID int64) ([]string, error)
}
