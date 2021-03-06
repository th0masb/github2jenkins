package g2j

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	JenkinsUrl   string `yaml:"jenkins-url"`
	Secrets      string
	Repositories []Repository `yaml:",flow"`
}

type Repository struct {
	Name     string
	Projects []Project `yaml:",flow"`
}

type Project struct {
	Path string
	Jobs []Job `yaml:",flow"`
}

type Job struct {
	Branch      string
	Name        string
	Parameters  string
	TokenKey    string `yaml:"token-key"`
	DiffMatcher string `yaml:"diff-matcher"`
}

func LoadConfig(path string) (Config, error) {
	fileBytes, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return Config{}, readErr
	} else {
		return interpretConfig(fileBytes)
	}
}

func interpretConfig(data []byte) (Config, error) {
	config := Config{}
	yamlErr := yaml.Unmarshal(data, &config)
	return config, yamlErr
}
