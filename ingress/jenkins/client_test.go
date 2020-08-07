package jenkins

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/th0masb/github2jenkins/tools"
)

var clientConfig = ClientConfig{
	JenkinsURL: "mock-url",
}

func TestSuccessfulNilParametersCall(t *testing.T) {
	// assemble
	clientRequest := TriggerJobRequest{
		JobName:    "mock-job",
		Cause:      "some-cause",
		Token:      "ABC",
		Parameters: nil,
	}

	expectedURL := "mock-url/job/mock-job/build?token=ABC&cause=some-cause"

	mockBody := tools.NewMockBody("")
	mockBody.On("Close").Return(nil).Once()

	mockRequester := tools.NewMockRequester()
	mockRequester.
		On("Do", createExpectedRequest(t, expectedURL)).
		Return(tools.Response(http.StatusOK, mockBody)).
		Once()

	underTest := Client{config: clientConfig, delegate: &mockRequester}

	// act
	err := underTest.TriggerJob(clientRequest)

	// assert
	assert.Nil(t, err)
	mockBody.AssertExpectations(t)
	mockRequester.AssertExpectations(t)
}

func TestSuccessfulEmptyParametersCall(t *testing.T) {
	// assemble
	clientRequest := TriggerJobRequest{
		JobName:    "mock-job",
		Cause:      "some-cause",
		Token:      "ABC",
		Parameters: map[string]string{},
	}

	expectedURL := "mock-url/job/mock-job/build?token=ABC&cause=some-cause"

	mockBody := tools.NewMockBody("")
	mockBody.On("Close").Return(nil).Once()

	mockRequester := tools.NewMockRequester()
	mockRequester.
		On("Do", createExpectedRequest(t, expectedURL)).
		Return(tools.Response(http.StatusOK, mockBody)).
		Once()

	underTest := Client{config: clientConfig, delegate: &mockRequester}

	// act
	err := underTest.TriggerJob(clientRequest)

	// assert
	assert.Nil(t, err)
	mockBody.AssertExpectations(t)
	mockRequester.AssertExpectations(t)
}

func TestSuccessfulParameterisedCall(t *testing.T) {
	// assemble
	clientRequest := TriggerJobRequest{
		JobName: "mock-job",
		Cause:   "some-cause",
		Token:   "ABC",
		Parameters: map[string]string{
			"param1": "value1",
			"param2": "value2",
		},
	}

	expectedURL := fmt.Sprintf(
		"%s%s",
		"mock-url/job/mock-job/buildWithParameters?token=ABC",
		"&cause=some-cause&param1=value1&param2=value2",
	)

	mockBody := tools.NewMockBody("")
	mockBody.On("Close").Return(nil).Once()

	mockRequester := tools.NewMockRequester()
	mockRequester.
		On("Do", createExpectedRequest(t, expectedURL)).
		Return(tools.Response(http.StatusOK, mockBody)).
		Once()

	underTest := Client{config: clientConfig, delegate: &mockRequester}

	// act
	err := underTest.TriggerJob(clientRequest)

	// assert
	assert.Nil(t, err)
	mockBody.AssertExpectations(t)
	mockRequester.AssertExpectations(t)
}

func createExpectedRequest(t *testing.T, url string) *http.Request {
	expectedRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.FailNow()
	}
	return expectedRequest
}
