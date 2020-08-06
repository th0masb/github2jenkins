package hook

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	hookTypeHeaderKey string = "X-Github-Event"
	pushHookID        string = "push"
	pullRequestID     string = "pull-request"
	pingHookID        string = "ping"
)

// Parse A hook request body sent from github
func Parse(headers http.Header, body []byte) (Hook, error) {
	switch reqType := headers.Get(hookTypeHeaderKey); reqType {
	case pushHookID:
		hook := Push{}
		err := json.Unmarshal(body, &hook)
		return hook, err
	case pingHookID:
		return Ping{}, nil
	default:
		return Push{}, fmt.Errorf("Unrecognised hook type: %s", reqType)
	}
}

// Push A push hook from github
type Push struct {
	Ref        string
	Before     string
	After      string
	Repository Repository
	Pusher     Pusher
}

// Pusher Information about the person who pushed any changes
type Pusher struct {
	Name  string
	Email string
}

// Repository Information about the repository that was pushed to
type Repository struct {
	Name  string
	Owner Owner
}

// Owner Information about the repository owner
type Owner struct {
	Name string
}

// Ping Represents a ping hook from github
type Ping struct {
}

// Hook Represents a hook from github
type Hook interface {
	isHook()
}

func (Push) isHook() {}
func (Ping) isHook() {}
