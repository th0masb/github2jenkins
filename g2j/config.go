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

// YamlConfig structure of config in yaml form
type YamlConfig struct {
	Jenkins      Jenkins
	Repositories []YamlRepository `yaml:",flow"`
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

// YamlRepository structure of repository in yaml form
type YamlRepository struct {
	Name     string
	Projects []YamlProject `yaml:",flow"`
}

// Project Configuration for a project contained within a repo
type Project struct {
	Path string
	Jobs []Job
}

// YamlProject structure of project in yaml form
type YamlProject struct {
	Path string
	Jobs []YamlJob `yaml:",flow"`
}

// Job Configuration for a Jenkins job triggered by changes to a project
type Job struct {
	Name          string
	Parameters    string
	Token         string
	BranchMatcher *regexp.Regexp
	DiffMatcher   *regexp.Regexp
}

// YamlJob structure of job in yaml form
type YamlJob struct {
	Name          string
	Parameters    string
	TokenKey      string `yaml:"token-key"`
	BranchMatcher string `yaml:"branch-matcher"`
	DiffMatcher   string `yaml:"diff-matcher"`
}

// LoadConfig Loads application configuration from a yaml file at the given
// location.
func LoadConfig(configPath, secretsPath string) (*Config, error) {
	fileBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return &Config{}, err
	}
	secrets, err := loadSecrets(secretsPath)
	if err != nil {
		return &Config{}, err
	}
	return interpretConfig(fileBytes, secrets)
}

func interpretConfig(data []byte, secrets Secrets) (*Config, error) {
	config := YamlConfig{}
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return &Config{}, err
	}
	//panic(fmt.Sprintf("%+v", config))
	return convertRawConfig(&config, secrets)
}

func convertRawConfig(rawConfig *YamlConfig, secrets Secrets) (*Config, error) {
	repositories, err := mapRepositorySlice(rawConfig.Repositories, secrets)
	if err != nil {
		return &Config{}, err
	}
	return &Config{
		Jenkins:      rawConfig.Jenkins,
		Secrets:      secrets,
		Repositories: repositories,
	}, nil
}

func mapRepositorySlice(src []YamlRepository, secrets Secrets) ([]Repository, error) {
	dest := make([]Repository, len(src))
	for i, r := range src {
		mappedProjects, err := mapProjectSlice(r.Projects, secrets)
		if err != nil {
			return dest, err
		}
		dest[i] = Repository{
			Name:     r.Name,
			Projects: mappedProjects,
		}
	}
	return dest, nil
}

func mapProjectSlice(src []YamlProject, secrets Secrets) ([]Project, error) {
	dest := make([]Project, len(src))
	for i, p := range src {
		mappedJobs, err := mapJobSlice(p.Jobs, secrets)
		if err != nil {
			return dest, err
		}
		dest[i] = Project{
			Path: p.Path,
			Jobs: mappedJobs,
		}
	}
	return dest, nil
}

func mapJobSlice(src []YamlJob, secrets Secrets) ([]Job, error) {
	dest := make([]Job, len(src))
	for i, j := range src {
		branchMatcher, err := regexp.Compile(j.BranchMatcher)
		if err != nil {
			return dest, err
		}
		diffMatcher, err := regexp.Compile(j.DiffMatcher)
		if err != nil {
			return dest, err
		}
		token, tokenWasPresent := secrets[j.TokenKey]
		if !tokenWasPresent {
			return dest, fmt.Errorf("Token with key %s not found for job %s", j.TokenKey, j.Name)
		}
		dest[i] = Job{
			Name:          j.Name,
			Parameters:    j.Parameters,
			Token:         token,
			BranchMatcher: branchMatcher,
			DiffMatcher:   diffMatcher,
		}
	}
	return dest, nil
}
