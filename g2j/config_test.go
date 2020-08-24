package g2j

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFailureDueToBadRegex(t *testing.T) {
	// assemble
	configYaml := `
        jenkins:
          url: x
        repositories:
         - name: guthub2jenkins
           projects:
             - path: first/
               jobs:
                 - branch-matcher: [
                   name: my-job
                   token-key: A
                   diff-matcher: ^src/*
    `

	secrets := map[string]string{
		"A": "TokenA",
	}

	// act
	config, err := interpretConfig([]byte(configYaml), secrets)

	// assert
	assert.NotNil(t, err)
	assert.Equal(t, &Config{}, config)
}

func TestParseFailureDueToMissingToken(t *testing.T) {
	// assemble
	configYaml := `
        jenkins:
          url: x
        repositories:
         - name: guthub2jenkins
           projects:
             - path: first/
               jobs:
                 - branch-matcher: master
                   name: my-job
                   token-key: A
                   diff-matcher: ^src/*
    `

	secrets := map[string]string{
		"B": "TokenA",
	}

	// act
	config, err := interpretConfig([]byte(configYaml), secrets)

	// assert
	assert.NotNil(t, err)
	assert.Equal(t, &Config{}, config)
}

func TestPartialConfig(t *testing.T) {
	// assemble
	configYaml := `
        jenkins:
          url: x
        repositories:
         - name: guthub2jenkins
           projects:
             - path: first/
               jobs:
                 - branch-matcher: master
                   name: my-job
                   token-key: A
                   diff-matcher: ^src/*
    `

	secrets := map[string]string{
		"A": "TokenA",
	}

	expectedConfig := &Config{
		Jenkins: Jenkins{
			URL:                "x",
			Protocol:           "",
			TLSCertificatePath: "",
		},
		Secrets: secrets,
		Repositories: []*Repository{
			&Repository{
				Name: "guthub2jenkins",
				Projects: []*Project{
					&Project{
						Path: "first/",
						Jobs: []*Job{
							&Job{
								BranchMatcher: regexp.MustCompile("master"),
								Name:          "my-job",
								Parameters:    "",
								Token:         "TokenA",
								DiffMatcher:   regexp.MustCompile(`^src/*`),
							},
						},
					},
				},
			},
		},
	}

	// act
	actualConfig, err := interpretConfig([]byte(configYaml), secrets)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedConfig, actualConfig)
}

func TestFullConfig(t *testing.T) {
	// assemble
	configYaml := `
        jenkins:
          url: https://myhost:8443
          protocol: https
          tls-cert: /path/to/cert
        repositories:
         - name: github2jenkins
           projects:
             - path: first/
               jobs:
                 - branch-matcher: master
                   name: my-job
                   parameters: Some expression
                   token-key: A
                   diff-matcher: src/.*
                 - branch-matcher: "master|dev"
                   name: job2
                   parameters: Other expression
                   token-key: B
                   diff-matcher: any
             - path: second/path/
               jobs:
                 - branch-matcher: master
                   name: my-job
                   parameters: Some expression
                   token-key: A
                   diff-matcher: src/ci/.*
         - name: github2jenkins2
           projects:
             - path: first/
               jobs:
                 - branch-matcher: master
                   name: my-job
                   parameters: Some expression
                   token-key: A
                   diff-matcher: src/*
                 - branch-matcher: ".*"
                   name: job2
                   parameters: Other expression
                   token-key: B
                   diff-matcher: any
    `

	secrets := map[string]string{
		"A": "TokenA",
		"B": "TokenB",
	}

	expectedConfig := &Config{
		Secrets: secrets,
		Jenkins: Jenkins{
			URL:                "https://myhost:8443",
			Protocol:           "https",
			TLSCertificatePath: "/path/to/cert",
		},
		Repositories: []*Repository{
			&Repository{
				Name: "github2jenkins",
				Projects: []*Project{
					&Project{
						Path: "first/",
						Jobs: []*Job{
							&Job{
								BranchMatcher: regexp.MustCompile("master"),
								Name:          "my-job",
								Parameters:    "Some expression",
								Token:         "TokenA",
								DiffMatcher:   regexp.MustCompile("src/.*"),
							},
							&Job{
								BranchMatcher: regexp.MustCompile("master|dev"),
								Name:          "job2",
								Parameters:    "Other expression",
								Token:         "TokenB",
								DiffMatcher:   regexp.MustCompile("any"),
							},
						},
					},
					&Project{
						Path: "second/path/",
						Jobs: []*Job{
							&Job{
								BranchMatcher: regexp.MustCompile("master"),
								Name:          "my-job",
								Parameters:    "Some expression",
								Token:         "TokenA",
								DiffMatcher:   regexp.MustCompile("src/ci/.*"),
							},
						},
					},
				},
			},
			&Repository{
				Name: "github2jenkins2",
				Projects: []*Project{
					&Project{
						Path: "first/",
						Jobs: []*Job{
							&Job{
								BranchMatcher: regexp.MustCompile("master"),
								Name:          "my-job",
								Parameters:    "Some expression",
								Token:         "TokenA",
								DiffMatcher:   regexp.MustCompile("src/*"),
							},
							&Job{
								BranchMatcher: regexp.MustCompile(".*"),
								Name:          "job2",
								Parameters:    "Other expression",
								Token:         "TokenB",
								DiffMatcher:   regexp.MustCompile("any"),
							},
						},
					},
				},
			},
		},
	}

	// act
	actualConfig, err := interpretConfig([]byte(configYaml), secrets)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedConfig, actualConfig)
}
