package mocks

import (
	"kube-job-runner/pkg/app/data"

	"github.com/stretchr/testify/mock"
)

type MockDataStore struct {
	mock.Mock
}

func (eventStore *MockDataStore) PutDocument(document data.Document) error {
	args := eventStore.Called(document)
	return args.Error(0)
}

func (eventStore *MockDataStore) GetDocuments(id string) ([]data.Document, error) {
	args := eventStore.Called(id)
	return args.Get(0).([]data.Document), args.Error(1)
}

func (eventStore *MockDataStore) GivenGetDocumentsReturns(documents []data.Document) {
	eventStore.On("GetDocuments", mock.Anything).Return(documents, nil)
}

func (eventStore *MockDataStore) GivenGetDocumentsByStatusReturns(documents []data.Document) {
	eventStore.On("GetDocumentsByStatus", mock.Anything).Return(documents)
}

func (eventStore *MockDataStore) GivenGetDocumentsFailed(err error) {
	eventStore.On("GetDocuments", mock.Anything).Return([]data.Document{}, err)
}

func (eventStore *MockDataStore) GivenPutDocumentFails(err error) {
	eventStore.On("PutDocument", mock.Anything).Return(err)
}

func (eventStore *MockDataStore) GivenPutDocumentSucceeds() {
	eventStore.On("PutDocument", mock.Anything).Return(nil)
}

func CreateMockStore() *MockDataStore {
	mockEventStore := new(MockDataStore)
	return mockEventStore
}
