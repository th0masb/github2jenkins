package tools

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockRequester wraps a mock for testing
type MockRequester struct{ mock.Mock }

// NewMockRequester creates a new mock requester
func NewMockRequester() MockRequester {
	return MockRequester{mock.Mock{}}
}

// Do notifies the mock that a call is made and retrieves the stubbed values
func (m *MockRequester) Do(request *http.Request) (*http.Response, error) {
	args := m.Called(request)
	return args.Get(0).(*http.Response), args.Error(1)
}
