package hook

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	hookTypeHeaderKey string = "X-Github-Event"
	pushHookId        string = "push"
	pullRequestId     string = "pull-request"
	pingHookId        string = "ping"
)

func Parse(headers http.Header, body []byte) (Hook, error) {
	switch reqType := headers.Get(hookTypeHeaderKey); reqType {
	case pushHookId:
		hook := Push{}
		err := json.Unmarshal(body, &hook)
		return hook, err
	case pingHookId:
		return Ping{}, nil
	default:
		return Push{}, fmt.Errorf("Unrecognised hook type: %s", reqType)
	}
}

type Push struct {
	Ref        string
	Before     string
	After      string
	Repository Repository
	Pusher     Pusher
}

type Pusher struct {
	Name  string
	Email string
}

type Repository struct {
	Name  string
	Owner Owner
}

type Owner struct {
	Name string
}

type Ping struct {
}

type Hook interface {
	isHook()
}

func (_ Push) isHook() {}
func (_ Ping) isHook() {}
