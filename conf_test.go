package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestPartialConfig(t *testing.T) {
	configYaml := `
        ci-dir: ci
        projects:
          - path: first/
            jobs:
              - branch: master
                name: my-job
                token: A
                diff-matcher: src/**
    `

	expectedConfig := Config{
		JenkinsUrl: "",
		DirCI:      "ci",
		Projects: []Project{
			Project{
				Path: "first/",
				Jobs: []Job{
					Job{
						Branch:      "master",
						Name:        "my-job",
						Parameters:  "",
						Token:       "A",
						DiffMatcher: "src/**",
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
        ci-dir: ci
        projects:
          - path: first/
            jobs:
              - branch: master
                name: my-job
                parameters: Some expression
                token: A
                diff-matcher: src/**
              - branch: "*"
                name: job2
                parameters: Other expression
                token: B
                diff-matcher: any
          - path: second/path/
            jobs:
              - branch: master
                name: my-job
                parameters: Some expression
                token: A
                diff-matcher: src/**
            `

	expectedConfig := Config{
		JenkinsUrl: "https://myhost:8443",
		DirCI:      "ci",
		Projects: []Project{
			Project{
				Path: "first/",
				Jobs: []Job{
					Job{
						Branch:      "master",
						Name:        "my-job",
						Parameters:  "Some expression",
						Token:       "A",
						DiffMatcher: "src/**",
					},
					Job{
						Branch:      "*",
						Name:        "job2",
						Parameters:  "Other expression",
						Token:       "B",
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
						Token:       "A",
						DiffMatcher: "src/**",
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
