package ingress

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/th0masb/github2jenkins/g2j"
	"github.com/th0masb/github2jenkins/ingress/diff"
	"github.com/th0masb/github2jenkins/ingress/hook"
)

// HookHandler handles github hook requests
type HookHandler struct {
	diffClient diffClientWrapper
	config     g2j.Config
	secrets    g2j.Secrets
}

type diffClientWrapper interface {
	RequestPushDiff(push *hook.Push) ([]string, error)
}

// NewHookHandler create a new github hook handler
func NewHookHandler(
	config g2j.Config,
	secrets g2j.Secrets,
) *HookHandler {
	return &HookHandler{
		diffClient: diff.CreateRestClient(),
		config:     config,
		secrets:    secrets,
	}
}

// Handle the incoming hook request from github
func (hh HookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Unable to read request body: %s\n", err)
		return
	}
	h, err := hook.Parse(r.Header, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Unable to parse request body: %s %s\n", err, body)
		return
	}
	switch v := h.(type) {
	case hook.Ping:
		log.Printf("Received ping hook\n")
		w.WriteHeader(http.StatusOK)
	case hook.Push:
		log.Printf("Received push hook, requesting diff\n")
		filesChanged, err := hh.diffClient.RequestPushDiff(&v)
		if err != nil {
			log.Printf("Error calling diff client: %s\n", err)
			w.WriteHeader(http.StatusFailedDependency)
		} else {
			log.Printf("Files changed: %s\n", filesChanged)
			w.WriteHeader(http.StatusOK)
		}
	}
}
