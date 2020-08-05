package g2j

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Config The configuration of the application
type Config struct {
	Jenkins      Jenkins
	Secrets      string
	Repositories []Repository `yaml:",flow"`
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
	Projects []Project `yaml:",flow"`
}

// Project Configuration for a project contained within a repo
type Project struct {
	Path string
	Jobs []Job `yaml:",flow"`
}

// Job Configuration for a Jenkins job triggered by changes to a project
type Job struct {
	Branch      string
	Name        string
	Parameters  string
	TokenKey    string `yaml:"token-key"`
	DiffMatcher string `yaml:"diff-matcher"`
}

// LoadConfig Loads application configuration from a yaml file at the given
// location.
func LoadConfig(path string) (Config, error) {
	fileBytes, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return Config{}, readErr
	}
	return interpretConfig(fileBytes)
}

func interpretConfig(data []byte) (Config, error) {
	config := Config{}
	yamlErr := yaml.Unmarshal(data, &config)
	return config, yamlErr
}
