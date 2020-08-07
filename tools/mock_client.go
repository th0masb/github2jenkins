package tools

import (
	"io"
	"net/http"
	"strings"

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
	args := m.MethodCalled("Do", request)
	return args.Get(0).(*http.Response), args.Error(1)
}

// MockBody An io.ReadCloser which allows to assert on closing
type MockBody struct {
	mock.Mock
	isClosed bool
	delegate io.Reader
}

func (mb *MockBody) Read(b []byte) (int, error) {
	if mb.isClosed {
		panic("Read after close")
	}
	return mb.delegate.Read(b)
}

func (mb *MockBody) Close() error {
	if mb.isClosed {
		panic("Double close")
	}
	mb.isClosed = true
	args := mb.MethodCalled("Close")
	return args.Error(0)
}

// NewMockBody Create a new body which allows assertions that closing has occurred
// exactly once after all reads.
func NewMockBody(data string) *MockBody {
	return &MockBody{
		mock.Mock{},
		false,
		strings.NewReader(data),
	}
}

// Response creates a mock http response with given status and body
func Response(statusCode int, body io.ReadCloser) (*http.Response, error) {
	resp := http.Response{
		Status:     string(statusCode),
		StatusCode: statusCode,
		Body:       body,
	}
	return &resp, nil
}
