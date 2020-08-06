package diff

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/th0masb/github2jenkins/ingress/hook"
	"github.com/th0masb/github2jenkins/tools"
)

const (
	before              string = "abc"
	after               string = "def"
	repoName            string = "repo"
	repoOwnerName       string = "th0masb"
	expectedAcceptValue string = "application/vnd.github.VERSION.diff"
	expectedBaseURL     string = "https://api.github.com/repos"
)

func TestRequestFailPath(t *testing.T) {
	// assemble
	pushHook := hook.Push{
		Before: before,
		After:  after,
		Repository: hook.Repository{
			Name: repoName,
			Owner: hook.Owner{
				Name: repoOwnerName,
			},
		},
	}

	mockRequester := tools.NewMockRequester()
	mockRequester.
		On("Do", request(&pushHook)).
		Return(response(http.StatusServiceUnavailable, "")).
		Once()

	underTest := Client{&mockRequester}

	// act
	changedFiles, err := underTest.RequestPushDiff(&pushHook)

	// assert
	assert.Nil(t, changedFiles)
	assert.NotNil(t, err)
	mockRequester.AssertExpectations(t)
}

func TestRequestHappyPath(t *testing.T) {
	// assemble
	expectedFiles := []string{"path/a/b", "second/b/path", "x/y/z", "h/j/k"}
	responseBody := fmt.Sprintf(
		"%s\n%s\n%s",
		"diff --git a/path/a/b b/second/b/path",
		"ignore this line",
		"diff --git a/x/y/z b/h/j/k",
	)
	pushHook := hook.Push{
		Before: before,
		After:  after,
		Repository: hook.Repository{
			Name: repoName,
			Owner: hook.Owner{
				Name: repoOwnerName,
			},
		},
	}

	mockRequester := tools.NewMockRequester()
	mockRequester.
		On("Do", request(&pushHook)).
		Return(response(http.StatusOK, responseBody)).
		Once()

	underTest := Client{&mockRequester}

	// act
	filesChanged, err := underTest.RequestPushDiff(&pushHook)

	// assert
	assert.Equal(t, expectedFiles, filesChanged)
	assert.Nil(t, err)
	mockRequester.AssertExpectations(t)
}

func request(push *hook.Push) *http.Request {
	url := fmt.Sprintf(
		"%s/%s/%s/compare/%s...%s",
		expectedBaseURL,
		push.Repository.Owner.Name,
		push.Repository.Name,
		push.Before, push.After,
	)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", expectedAcceptValue)
	return req
}

func response(statusCode int, body string) (*http.Response, error) {
	resp := http.Response{
		Status:     string(statusCode),
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
	}
	return &resp, nil
}
