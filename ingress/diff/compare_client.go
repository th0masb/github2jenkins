package diff

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/th0masb/github2jenkins/ingress/hook"
	"github.com/th0masb/github2jenkins/tools"
)

const (
	diffHeaderKey      string = "Accept"
	diffHeaderValue    string = "application/vnd.github.VERSION.diff"
	baseURL            string = "https://api.github.com/repos"
	getMethod          string = "GET"
	diffRegexPattern   string = `^diff --git.+$`
	changedFilePattern string = `\s[ab](/\S+)+`
)

var diffRegex = regexp.MustCompile(diffRegexPattern)
var changedFileRegex = regexp.MustCompile(changedFilePattern)

// Client Testable wrapper around http requests
type Client struct{ tools.HTTPRequester }

// CreateRestClient Creates a rest client to make http calls
func CreateRestClient() *Client {
	return &Client{tools.NewRestRequester()}
}

// RequestPushDiff Fetches diff caused by github push
func (c *Client) RequestPushDiff(push *hook.Push) ([]string, error) {
	repoName := push.Repository.Name
	ownerName := push.Repository.Owner.Name
	before, after := push.Before, push.After
	url := fmt.Sprintf(
		"%s/%s/%s/compare/%s...%s",
		baseURL, ownerName, repoName, before, after,
	)
	body, err := c.getDiff(url)
	if err != nil {
		return nil, err
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

func (c *Client) getDiff(url string) (string, error) {
	req, err := http.NewRequest(getMethod, url, nil)
	req.Header.Add(diffHeaderKey, diffHeaderValue)
	if err != nil {
		return "", err
	}
	resp, err := c.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	if tools.IsErrorCode(resp.StatusCode) {
		log.Printf("Bad response: %s\n", resp.Status)
		return "", fmt.Errorf("Bad response: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}
