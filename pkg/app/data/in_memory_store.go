package data

import (
	"fmt"
	"sort"
	"sync"
)

type InMemoryStore struct {
	documentsList []Document
	mux           sync.Mutex
}

func (store *InMemoryStore) PutDocument(document Document) error {
	store.mux.Lock()
	defer store.mux.Unlock()
	store.documentsList = append(store.documentsList, document)
	return nil
}

func (store *InMemoryStore) GetDocuments(id string) ([]Document, error) {
	documents := make([]Document, 0)
	for _, document := range store.documentsList {
		if document.JobID == id {
			documents = append(documents, document)
		}
	}

	if len(documents) == 0 {
		return []Document{}, fmt.Errorf("no document found with id '%v'", id)
	}

	sort.Slice(documents, func(i, j int) bool {
		return documents[i].Timestamp > documents[j].Timestamp
	})

	return documents, nil
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		documentsList: []Document{},
	}
}
