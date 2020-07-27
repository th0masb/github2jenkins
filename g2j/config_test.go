package g2j

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestPartialConfig(t *testing.T) {
	configYaml := `
        secrets: "/path/to/secrets.env"
        repositories:
         - name: guthub2jenkins
           projects:
             - path: first/
               jobs:
                 - branch: master
                   name: my-job
                   token-key: A
                   diff-matcher: src/**
    `

	expectedConfig := Config{
		JenkinsUrl: "",
		Secrets:    "/path/to/secrets.env",
		Repositories: []Repository{
			Repository{
				Name: "guthub2jenkins",
				Projects: []Project{
					Project{
						Path: "first/",
						Jobs: []Job{
							Job{
								Branch:      "master",
								Name:        "my-job",
								Parameters:  "",
								TokenKey:    "A",
								DiffMatcher: "src/**",
							},
						},
					},
				},
			},
		},
	}

	assertParsedConfigAsExpected(&expectedConfig, configYaml, t)

}

func TestFullConfig(t *testing.T) {
	configYaml := `
        jenkins-url: https://myhost:8443
        secrets: /path/to/secrets.env
        repositories:
         - name: github2jenkins
           projects:
             - path: first/
               jobs:
                 - branch: master
                   name: my-job
                   parameters: Some expression
                   token-key: A
                   diff-matcher: src/**
                 - branch: "*"
                   name: job2
                   parameters: Other expression
                   token-key: B
                   diff-matcher: any
             - path: second/path/
               jobs:
                 - branch: master
                   name: my-job
                   parameters: Some expression
                   token-key: A
                   diff-matcher: src/**
         - name: github2jenkins2
           projects:
             - path: first/
               jobs:
                 - branch: master
                   name: my-job
                   parameters: Some expression
                   token-key: A
                   diff-matcher: src/**
                 - branch: "*"
                   name: job2
                   parameters: Other expression
                   token-key: B
                   diff-matcher: any
            `

	expectedConfig := Config{
		JenkinsUrl: "https://myhost:8443",
		Secrets:    "/path/to/secrets.env",
		Repositories: []Repository{
			Repository{
				Name: "github2jenkins",
				Projects: []Project{
					Project{
						Path: "first/",
						Jobs: []Job{
							Job{
								Branch:      "master",
								Name:        "my-job",
								Parameters:  "Some expression",
								TokenKey:    "A",
								DiffMatcher: "src/**",
							},
							Job{
								Branch:      "*",
								Name:        "job2",
								Parameters:  "Other expression",
								TokenKey:    "B",
								DiffMatcher: "any",
							},
						},
					},
					Project{
						Path: "second/path/",
						Jobs: []Job{
							Job{
								Branch:      "master",
								Name:        "my-job",
								Parameters:  "Some expression",
								TokenKey:    "A",
								DiffMatcher: "src/**",
							},
						},
					},
				},
			},
			Repository{
				Name: "github2jenkins2",
				Projects: []Project{
					Project{
						Path: "first/",
						Jobs: []Job{
							Job{
								Branch:      "master",
								Name:        "my-job",
								Parameters:  "Some expression",
								TokenKey:    "A",
								DiffMatcher: "src/**",
							},
							Job{
								Branch:      "*",
								Name:        "job2",
								Parameters:  "Other expression",
								TokenKey:    "B",
								DiffMatcher: "any",
							},
						},
					},
				},
			},
		},
	}

	assertParsedConfigAsExpected(&expectedConfig, configYaml, t)
}

func assertParsedConfigAsExpected(expected *Config, yaml string, t *testing.T) {
	actual, err := interpretConfig([]byte(yaml))
	if err != nil {
		t.Errorf("Error during parsing: %s\n", err)
	} else if !cmp.Equal(*expected, actual) {
		t.Errorf("Expected:\n%s\nbut received:\n%s\n", *expected, actual)
	}
}
