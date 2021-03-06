package diff

import (
	"fmt"
	"github.com/th0masb/github2jenkins/hook"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	diffHeaderKey      string = "Accept"
	diffHeaderValue    string = "application/vnd.github.VERSION.diff"
	baseUrl            string = "https://api.github.com/repos"
	getMethod          string = "GET"
	diffRegexPattern   string = `^diff --git.+$`
	changedFilePattern string = `\s[ab](/\S+)+`
)

var diffRegex = regexp.MustCompile(diffRegexPattern)
var changedFileRegex = regexp.MustCompile(changedFilePattern)

type Client struct{ requester Requester }

type Requester interface {
	Get(request *http.Request) (*http.Response, error)
}

type requesterImpl struct{ delegate *http.Client }

func CreateRestClient() Client {
	httpClient := http.Client{}
	return Client{
		requester: requesterImpl{
			delegate: &httpClient,
		},
	}
}

func (c Client) RequestPushDiff(push *hook.Push) ([]string, error) {
	repoName := push.Repository.Name
	ownerName := push.Repository.Owner.Name
	before, after := push.Before, push.After
	url := fmt.Sprintf("%s/%s/%s/compare/%s...%s", baseUrl, ownerName, repoName, before, after)
	body, err := c.getDiff(url)
	if err != nil {
		return make([]string, 0), err
	}
	lines := strings.Split(string(body), "\n")
	return extractChangedFiles(lines), err
}

func extractChangedFiles(diffLines []string) []string {
	changedFiles := make([]string, 0)
	uniqueChangedFiles := make(map[string]bool)
	for _, line := range diffLines {
		if diffRegex.Match([]byte(line)) {
			for _, changedFile := range changedFileRegex.FindAll([]byte(line), -1) {
				changedFile := string(changedFile)[3:]
				if !uniqueChangedFiles[changedFile] {
					uniqueChangedFiles[changedFile] = true
					changedFiles = append(changedFiles, changedFile)
				}
			}
		}
	}
	return changedFiles
}

func (c Client) getDiff(url string) (string, error) {
	req, err := http.NewRequest(getMethod, url, nil)
	req.Header.Add(diffHeaderKey, diffHeaderValue)
	if err != nil {
		return "", err
	}
	resp, err := c.requester.Get(req)
	if err != nil {
		return "", err
	}
	if isErrorCode(resp.StatusCode) {
		log.Printf("Bad response: %s\n", resp.Status)
		return "", fmt.Errorf("Bad response: %s\n", resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func isErrorCode(code int) bool {
	return code < 200 || code > 299
}

func (ri requesterImpl) Get(request *http.Request) (*http.Response, error) {
	return ri.delegate.Do(request)
}
