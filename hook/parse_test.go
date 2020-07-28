package hook

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
	hook, err := Parse(header, body)

	// assert
	if err != nil {
		t.Fatalf("Expected success: %s\n", err)
	}

	if hook != (Ping{}) {
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
                "name": "github2jenkins",
				"owner": {
					"name": "th0masb"
				}
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
	hook, err := Parse(header, body)

	// assert
	if err != nil {
		t.Fatalf("Expected success: %s\n", err)
	}

	expected := Push{
		Ref:    "123",
		Before: "beforeHash",
		After:  "afterHash",
		Repository: Repository{
			Name: "github2jenkins",
			Owner: Owner{
				Name: "th0masb",
			},
		},
		Pusher: Pusher{
			Name:  "Tom",
			Email: "Email",
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
