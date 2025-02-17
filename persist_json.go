package go_memkvstore

import (
	"encoding/json"
	"os"
)

type JSONStorePersister[T any] struct {
	filename string
	store    *Store[T]
}

func NewJSONStorePersister[T any](filename string, kvstore *Store[T]) *JSONStorePersister[T] {
	return &JSONStorePersister[T]{filename: filename, store: kvstore}
}

func (m *JSONStorePersister[V]) Write() error {
	// open file with create, write and truncate flags
	f, err := os.OpenFile(m.filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	// close and sync file after function ends
	defer f.Close()
	defer f.Sync()
	// lock the store
	m.store.RLock()
	defer m.store.RUnlock()
	// create a new json encoder with the file as destination and encode the store
	enc := json.NewEncoder(f)
	// format the json output
	enc.SetIndent("", "  ")
	return enc.Encode(m.store.Store)
}

func (m *JSONStorePersister[V]) Read() error {
	// open file with read flag
	f, err := os.OpenFile(m.filename, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	// close file after function ends
	defer f.Close()
	// lock the store
	m.store.Lock()
	defer m.store.Unlock()
	// create a new json decoder with the file as source and decode the store
	dec := json.NewDecoder(f)
	return dec.Decode(&m.store.Store)
}
