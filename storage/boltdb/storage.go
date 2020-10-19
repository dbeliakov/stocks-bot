package boltdb

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

const (
	StatesBucketName  = "states"
	SymbolsBucketName = "symbols"
)

type Storage struct {
	db *bolt.DB
}

func NewStorage(path string) (*Storage, error) {
	db, err := bolt.Open(path, 0644, bolt.DefaultOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Init() error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(StatesBucketName)); err != nil {
			return fmt.Errorf("failed to create %s bucket: %w", StatesBucketName, err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(SymbolsBucketName)); err != nil {
			return fmt.Errorf("failed to create %s bucket: %w", SymbolsBucketName, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}
	return nil
}

func (s *Storage) SetState(chatID int64, state string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket([]byte(StatesBucketName)).Put(
			chatIDKey(chatID), []byte(state)); err != nil {

			return fmt.Errorf("failed to update value: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update state: %w", err)
	}
	return nil
}

func (s *Storage) GetState(chatID int64) (string, error) {
	var state string
	err := s.db.View(func(tx *bolt.Tx) error {
		state = string(tx.Bucket([]byte(StatesBucketName)).Get(chatIDKey(chatID)))
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to get state: %w", err)
	}
	return state, nil
}

func (s *Storage) AddSymbol(chatID int64, symbol string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte(SymbolsBucketName)).Get(chatIDKey(chatID))
		var list []string
		if len(data) > 0 {
			if err := json.Unmarshal(data, &list); err != nil {
				return fmt.Errorf("invalid data: %w", err)
			}
		}
		for _, s := range list {
			if s == symbol {
				return nil
			}
		}
		list = append(list, symbol)
		data, err := json.Marshal(list)
		if err != nil {
			return fmt.Errorf("failed to marshal json: %w", err)
		}
		if err := tx.Bucket([]byte(SymbolsBucketName)).Put(chatIDKey(chatID), data); err != nil {
			return fmt.Errorf("failed to update db value: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to add symbol: %w", err)
	}
	return nil
}

func (s *Storage) RemoveSymbol(chatID int64, symbol string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte(SymbolsBucketName)).Get(chatIDKey(chatID))
		var list []string
		if len(data) > 0 {
			if err := json.Unmarshal(data, &list); err != nil {
				return fmt.Errorf("invalid data: %w", err)
			}
		}
		var removed = false
		for i, s := range list {
			if s == symbol {
				list = append(list[:i], list[i+1:]...)
				removed = true
			}
		}
		if !removed {
			return nil
		}

		data, err := json.Marshal(list)
		if err != nil {
			return fmt.Errorf("failed to marshal json: %w", err)
		}
		if err := tx.Bucket([]byte(SymbolsBucketName)).Put(chatIDKey(chatID), data); err != nil {
			return fmt.Errorf("failed to update db value: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to add symbol: %w", err)
	}
	return nil
}

func (s *Storage) Symbols(chatID int64) ([]string, error) {
	var list []string
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte(SymbolsBucketName)).Get(chatIDKey(chatID))
		if len(data) == 0 {
			return nil
		}
		if err := json.Unmarshal(data, &list); err != nil {
			return fmt.Errorf("invalid data: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get symbols: %w", err)
	}
	return list, nil
}

func chatIDKey(id int64) []byte {
	return []byte(fmt.Sprintf("%d", id))
}
