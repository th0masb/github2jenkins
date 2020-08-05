package g2j

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartialConfig(t *testing.T) {
	// assemble
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
		Jenkins: Jenkins{
			URL:                "",
			Protocol:           "",
			TLSCertificatePath: "",
		},
		Secrets: "/path/to/secrets.env",
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

	// act
	actualConfig, err := interpretConfig([]byte(configYaml))

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedConfig, actualConfig)
}

func TestFullConfig(t *testing.T) {
	// assemble
	configYaml := `
        secrets: "/path/to/secrets.json"
        jenkins: 
          url: https://myhost:8443
          protocol: https
          tls-cert: /path/to/cert
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
		Secrets: "/path/to/secrets.json",
		Jenkins: Jenkins{
			URL:                "https://myhost:8443",
			Protocol:           "https",
			TLSCertificatePath: "/path/to/cert",
		},
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

	// act
	actualConfig, err := interpretConfig([]byte(configYaml))

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedConfig, actualConfig)
}
