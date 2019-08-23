package mocks

import "github.com/stretchr/testify/mock"

type MockIDGenerator struct {
	mock.Mock
}

func (generator *MockIDGenerator) Generate() string {
	args := generator.Called()
	return args.String(0)
}

func (generator *MockIDGenerator) GivenGeneratorReturnsID(id string) {
	generator.ExpectedCalls = []*mock.Call{}
	generator.On("Generate").Return(id)
}

func CreateMockIDGenerator() *MockIDGenerator {
	mockIdGenerator := new(MockIDGenerator)
	mockIdGenerator.GivenGeneratorReturnsID("default-id")
	return mockIdGenerator
}
