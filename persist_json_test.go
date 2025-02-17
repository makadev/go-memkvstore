package go_memkvstore_test

import (
	"os"
	"testing"
	"time"

	mkvstore "github.com/makadev/go-memkvstore"
)

func TestPersistStore(t *testing.T) {
	// create a store and write a value
	mkv := mkvstore.New[string](1 * time.Minute)
	mkv.Set("test", "value")

	value, ok := mkv.Get("test", "")
	if !ok {
		t.Errorf("Expected value to be found")
	}
	if value != "value" {
		t.Errorf("Expected value to be 'value'")
	}

	// persist the store
	persist := mkvstore.NewJSONStorePersister("test.json", mkv)
	defer os.Remove("test.json")

	err := persist.Write()
	if err != nil {
		t.Errorf("Expected write to succeed")
	}

	// create a new store and read the value
	mkv2 := mkvstore.New[string](1 * time.Minute)
	persist2 := mkvstore.NewJSONStorePersister("test.json", mkv2)
	err = persist2.Read()
	if err != nil {
		t.Errorf("Expected read to succeed")
	}

	value, ok = mkv2.Get("test", "")
	if !ok {
		t.Errorf("Expected value to be found")
	}
	if value != "value" {
		t.Errorf("Expected value to be 'value'")
	}
}

func TestPersistStoreAgain(t *testing.T) {
	// create a store and write a value
	mkv := mkvstore.New[string](1 * time.Minute)
	mkv.Set("test", "value")

	value, ok := mkv.Get("test", "")
	if !ok {
		t.Errorf("Expected value to be found")
	}
	if value != "value" {
		t.Errorf("Expected value to be 'value'")
	}

	// persist the store
	persist := mkvstore.NewJSONStorePersister("test.json", mkv)
	defer os.Remove("test.json")

	err := persist.Write()
	if err != nil {
		t.Errorf("Expected write to succeed")
	}

	// change value
	mkv.Set("test", "value2")
	// persist the store again, this should overwrite the file
	persist2 := mkvstore.NewJSONStorePersister("test.json", mkv)
	err = persist2.Write()
	if err != nil {
		t.Errorf("Expected write to succeed")
	}

	// create a new store and read the value
	mkv3 := mkvstore.New[string](1 * time.Minute)
	persist3 := mkvstore.NewJSONStorePersister("test.json", mkv3)
	err = persist3.Read()
	if err != nil {
		t.Errorf("Expected read to succeed")
	}

	value, ok = mkv3.Get("test", "")
	if !ok {
		t.Errorf("Expected value to be found")
	}
	if value != "value2" {
		t.Errorf("Expected value to be 'value2'")
	}
}
