package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIndexerAdd(t *testing.T) {
	indices := []Index{
		UniqueIndex{
			key:        "email",
			typeName:   "User",
			rootFolder: "accounts",
		},
		UniqueIndex{
			key:        "username",
			typeName:   "User",
			rootFolder: "accounts",
		},
		UniqueIndex{
			key:        "price",
			typeName:   "Product",
			rootFolder: "products",
		},
	}

	indexer := NewIndexer("/tmp/foo")

	for k := range indices {
		indexer.AddIndex(indices[k])
	}

	idxMap := indexer.indices

	exp, got := 2, len(idxMap)
	if exp != got {
		t.Fatalf("Expected %v different indices types for type User got %v", exp, got)
	}

	exp, got = 2, len(idxMap["User"])
	if exp != got {
		t.Fatalf("Expected %v indices for type User got %v", exp, got)
	}

	assert.ElementsMatch(t, indices, []Index{
		idxMap["User"]["email"][0],
		idxMap["User"]["username"][0],
		idxMap["Product"]["price"][0],
	})
}
