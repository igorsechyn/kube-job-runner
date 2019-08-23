package mocks

type AllMocks struct {
	MockJobClient   *MockJobClient
	MockStore       *MockDataStore
	MockLoggerSink  *MockLoggerSink
	MockQueueClient *MockQueueClient
	MockClock       *MockClock
	MockIdGenerator *MockIDGenerator
}

func InitMocks() AllMocks {
	mockJobClient := CreateMockJobClient()
	mockStore := CreateMockStore()
	mockLoggerSink := CreateMockLoggerSink()
	mockQueueClient := CreateMockQueueClient()
	mockClock := CreateMockClock()
	mockIdGenerator := CreateMockIDGenerator()

	return AllMocks{
		MockClock:       mockClock,
		MockJobClient:   mockJobClient,
		MockLoggerSink:  mockLoggerSink,
		MockQueueClient: mockQueueClient,
		MockStore:       mockStore,
		MockIdGenerator: mockIdGenerator,
	}
}
