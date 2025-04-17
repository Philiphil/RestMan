package gormrepository_test

import (
	"testing"

	"github.com/philiphil/restman/orm/gormrepository"
)

func TestChunkSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	chunkSize := 3

	chunks := gormrepository.ChunkSlice(slice, chunkSize)
	if len(chunks) != 4 {
		t.Errorf("Expected 4 chunks, got %d", len(chunks))
	}

	if len(chunks[0]) != 3 {
		t.Errorf("Expected 3 elements in first chunk, got %d", len(chunks[0]))
	}

	if len(chunks[3]) != 1 {
		t.Errorf("Expected 1 element in last chunk, got %d", len(chunks[3]))
	}
}
