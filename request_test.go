package main

import (
	"github.com/google/go-cmp/cmp"
	"net/http"
	"testing"
)

func TestPingHookParse(t *testing.T) {
	// assemble
	header := createHeader("X-Github-Event", "ping")
	body := []byte("")

	// action
	hook, err := ParseRequest(header, body)

	// assert
	if err != nil {
		t.Fatalf("Expected success: %s\n", err)
	}

	if hook != (PingHook{}) {
		t.Fatalf("Expected empty ping hook struct: %+v\n", hook)
	}
}

func TestPushHookParse(t *testing.T) {
	// assemble
	header := createHeader("X-Github-Event", "push")
	body := []byte(
		`
        {
            "unknown": "x",
            "ref": "123",
            "before": "beforeHash",
            "after": "afterHash",
            "repository": {
                "name": "github2jenkins"
            },
            "pusher": {
                "name": "Tom",
                "email": "Email"
            },
            "compare": "compare-url",
            "commits": [
            {
                "id": "1",
                "added": [],
                "removed": [],
                "modified": ["first"]
            },
            {
                "id": "2",
                "added": ["something"],
                "removed": [],
                "modified": []
            }
            ]
        }
        `,
	)

	// action
	hook, err := ParseRequest(header, body)

	// assert
	if err != nil {
		t.Fatalf("Expected success: %s\n", err)
	}

	expected := PushHook{
		Ref:    "123",
		Before: "beforeHash",
		After:  "afterHash",
		Repository: Repository{
			Name: "github2jenkins",
		},
		Pusher: Pusher{
			Name:  "Tom",
			Email: "Email",
		},
		Compare: "compare-url",
		Commits: []Commit{
			Commit{
				Id:       "1",
				Added:    []string{},
				Removed:  []string{},
				Modified: []string{"first"},
			},
			Commit{
				Id:       "2",
				Added:    []string{"something"},
				Removed:  []string{},
				Modified: []string{},
			},
		},
	}

	if !cmp.Equal(expected, hook) {
		t.Fatalf("Expected:\n%+v\nbut received:\n%+v\n", expected, hook)
	}
}

func createHeader(key, value string) http.Header {
	var header http.Header
	header = make(map[string][]string)
	header.Add(key, value)
	return header
}
