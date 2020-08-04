package diff

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/th0masb/github2jenkins/diff/mock_diff"
	"github.com/th0masb/github2jenkins/hook"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const (
	before              string = "abc"
	after               string = "def"
	repoName            string = "repo"
	repoOwnerName       string = "th0masb"
	expectedAcceptValue string = "application/vnd.github.VERSION.diff"
	expectedBaseUrl     string = "https://api.github.com/repos"
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRequester := mock_diff.NewMockRequester(ctrl)
	mockRequester.
		EXPECT().
		Get(gomock.Eq(request(&pushHook))).
		Return(response(http.StatusServiceUnavailable, ""))

	underTest := Client{requester: mockRequester}

	// act
	_, err := underTest.RequestPushDiff(&pushHook)

	// assert
	if err == nil {
		t.Errorf("Expected error but received nil\n")
	}
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRequester := mock_diff.NewMockRequester(ctrl)
	mockRequester.
		EXPECT().
		Get(gomock.Eq(request(&pushHook))).
		Return(response(http.StatusOK, responseBody))

	underTest := Client{requester: mockRequester}

	// act
	filesChanged, _ := underTest.RequestPushDiff(&pushHook)

	// assert
	if !cmp.Equal(filesChanged, expectedFiles) {
		t.Errorf("Received: %s\n", filesChanged)
	}
}

func request(push *hook.Push) *http.Request {
	url := fmt.Sprintf(
		"%s/%s/%s/compare/%s...%s",
		expectedBaseUrl,
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
