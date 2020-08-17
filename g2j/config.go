package g2j

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v3"
)

// Config The configuration of the application
type Config struct {
	Jenkins      Jenkins
	Secrets      Secrets
	Repositories []Repository
}

type yamlConfig struct {
	jenkins      Jenkins
	secrets      string
	repositories []yamlRepository `yaml:",flow"`
}

// Jenkins Configuration for Jenkins communication
type Jenkins struct {
	URL                string
	Protocol           string
	TLSCertificatePath string `yaml:"tls-cert"`
}

// Repository Configuration for a specific repository
type Repository struct {
	Name     string
	Projects []Project
}

type yamlRepository struct {
	name     string
	projects []yamlProject `yaml:",flow"`
}

// Project Configuration for a project contained within a repo
type Project struct {
	Path string
	Jobs []Job
}

type yamlProject struct {
	path string
	jobs []yamlJob `yaml:",flow"`
}

// Job Configuration for a Jenkins job triggered by changes to a project
type Job struct {
	Name          string
	Parameters    string
	Token         string
	BranchMatcher *regexp.Regexp
	DiffMatcher   *regexp.Regexp
}

type yamlJob struct {
	name          string
	parameters    string
	tokenKey      string `yaml:"token-key"`
	branchMatcher string `yaml:"branch-matcher"`
	diffMatcher   string `yaml:"diff-matcher"`
}

// LoadConfig Loads application configuration from a yaml file at the given
// location.
func LoadConfig(path string) (*Config, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return &Config{}, err
	}
	return interpretConfig(fileBytes)
}

func interpretConfig(data []byte) (*Config, error) {
	config := yamlConfig{}
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return &Config{}, err
	}
	secrets, err := loadSecrets(config.secrets)
	if err != nil {
		return &Config{}, err
	}
	return convertRawConfig(&config, secrets)
}

func convertRawConfig(rawConfig *yamlConfig, secrets Secrets) (*Config, error) {
	repositories, err := mapRepositorySlice(rawConfig.repositories, secrets)
	if err != nil {
		return &Config{}, err
	}
	return &Config{
		Jenkins:      rawConfig.jenkins,
		Secrets:      secrets,
		Repositories: repositories,
	}, nil
}

func mapRepositorySlice(src []yamlRepository, secrets Secrets) ([]Repository, error) {
	dest := make([]Repository, len(src))
	for i, r := range src {
		mappedProjects, err := mapProjectSlice(r.projects, secrets)
		if err != nil {
			return dest, err
		}
		dest[i] = Repository{
			Name:     r.name,
			Projects: mappedProjects,
		}
	}
	return dest, nil
}

func mapProjectSlice(src []yamlProject, secrets Secrets) ([]Project, error) {
	dest := make([]Project, len(src))
	for i, p := range src {
		mappedJobs, err := mapJobSlice(p.jobs, secrets)
		if err != nil {
			return dest, err
		}
		dest[i] = Project{
			Path: p.path,
			Jobs: mappedJobs,
		}
	}
	return dest, nil
}

func mapJobSlice(src []yamlJob, secrets Secrets) ([]Job, error) {
	dest := make([]Job, len(src))
	for i, j := range src {
		branchMatcher, err := regexp.Compile(j.branchMatcher)
		if err != nil {
			return dest, err
		}
		diffMatcher, err := regexp.Compile(j.diffMatcher)
		if err != nil {
			return dest, err
		}
		token, tokenWasPresent := secrets[j.tokenKey]
		if !tokenWasPresent {
			return dest, fmt.Errorf("Token with key %s not found for job %s", j.tokenKey, j.name)
		}
		dest[i] = Job{
			Name:          j.name,
			Parameters:    j.parameters,
			Token:         token,
			BranchMatcher: branchMatcher,
			DiffMatcher:   diffMatcher,
		}
	}
	return dest, nil
}
