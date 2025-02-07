package go_memkvstore_test

import (
	"fmt"
	"testing"
	"time"

	mkvstore "github.com/makadev/go-memkvstore"
)

func TestStore(t *testing.T) {
	mkv := mkvstore.New[string](1 * time.Minute)
	mkv.Set("test", "value")

	value, ok := mkv.Get("test", "")
	if !ok {
		t.Errorf("Expected value to be found")
	}
	if value != "value" {
		t.Errorf("Expected value to be 'value'")
	}
}

func TestMKVSExpiration(t *testing.T) {
	mkv := mkvstore.New[string](100 * time.Millisecond)
	mkv.Set("test", "value")

	value, ok := mkv.Get("test", "")
	if !ok {
		t.Errorf("Expected value to be found")
	}
	if value != "value" {
		t.Errorf("Expected value to be 'value'")
	}

	// Wait for expiration
	time.Sleep(200 * time.Millisecond)

	value, ok = mkv.Get("test", "")
	if ok {
		t.Errorf("Expected value to be expired")
	}
	if value != "" {
		t.Errorf("Expected value to be empty")
	}
}

func TestMKVSDelete(t *testing.T) {
	mkv := mkvstore.New[string](1 * time.Minute)
	mkv.Set("test", "value")

	value, ok := mkv.Get("test", "")
	if !ok {
		t.Errorf("Expected value to be found")
	}
	if value != "value" {
		t.Errorf("Expected value to be 'value'")
	}

	mkv.Delete("test")

	value, ok = mkv.Get("test", "")
	if ok {
		t.Errorf("Expected value to be deleted")
	}
	if value != "" {
		t.Errorf("Expected value to be empty")
	}
}

func TestMKVSCleanup(t *testing.T) {
	mkv := mkvstore.New[string](1 * time.Millisecond)
	mkv.Set("test", "value")
	mkv.SetWithExpiration("test2", "value2", 2000*time.Millisecond)

	value, ok := mkv.Get("test", "")
	if !ok {
		t.Errorf("Expected value to be found")
	}
	if value != "value" {
		t.Errorf("Expected value to be 'value'")
	}

	value, ok = mkv.Get("test2", "")
	if !ok {
		t.Errorf("Expected value to be found")
	}
	if value != "value2" {
		t.Errorf("Expected value to be 'value2'")
	}

	// Wait for expiration
	time.Sleep(200 * time.Millisecond)

	mkv.Cleanup()

	// test should be expired and deleted
	value, ok = mkv.Get("test", "")
	if ok {
		t.Errorf("Expected value to be deleted")
	}
	if value != "" {
		t.Errorf("Expected value to be empty")
	}

	// test2 should still be there
	value, ok = mkv.Get("test2", "")
	if !ok {
		t.Errorf("Expected value to be found")
	}
	if value != "value2" {
		t.Errorf("Expected value to be 'value2'")
	}
}

func BenchmarkAddMKVS(b *testing.B) {
	mkv := mkvstore.New[string](1 * time.Hour)
	testdata := make([]string, 1000)

	for i := 0; i < 1000; i++ {
		testdata[i] = "test" + fmt.Sprintf("%d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mkv.Set(testdata[i%1000], "value")
	}
}

func BenchmarkLookupMKVS(b *testing.B) {
	mkv := mkvstore.New[string](1 * time.Hour)
	testdata := make([]string, 1000)

	for idx := range testdata {
		testdata[idx] = "test" + fmt.Sprintf("%d", idx)
		mkv.Set(testdata[idx], "value")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mkv.Get(testdata[i%1000], "value")
	}
}

func BenchmarkDeletesMKVS(b *testing.B) {
	mkv := mkvstore.New[string](1 * time.Hour)
	testdata := make([]string, 1000)

	for idx := range testdata {
		testdata[idx] = "test" + fmt.Sprintf("%d", idx)
		mkv.Set(testdata[idx], "")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mkv.Delete(testdata[i%1000])
	}
}
