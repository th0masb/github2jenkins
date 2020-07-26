package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const hookTypeHeaderKey string = "X-Github-Event"
const pushHookId string = "push"
const pullRequestId string = "pull-request"
const pingHookId string = "ping"

func ParseRequest(headers http.Header, body []byte) (Hook, error) {
	switch reqType := headers.Get(hookTypeHeaderKey); reqType {
	case pushHookId:
		hook := PushHook{}
		err := json.Unmarshal(body, &hook)
		return hook, err
	case pingHookId:
		return PingHook{}, nil
	default:
		return PushHook{}, fmt.Errorf("Unrecognised hook type: %s", reqType)
	}
}

type PushHook struct {
	Ref        string     `json:"ref"`
	Before     string     `json:"before"`
	After      string     `json:"after"`
	Repository Repository `json:"repository"`
	Pusher     Pusher     `json:"pusher"`
	Compare    string     `json:"compare"`
	Commits    []Commit   `json:"commits"`
}

type Commit struct {
	Id       string   `json:"id"`
	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`
}

type Pusher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Repository struct {
	Name string `json:"name"`
}

type PingHook struct {
}

type Hook interface {
	isHook()
}

func (_ PushHook) isHook() {}
func (_ PingHook) isHook() {}
