package conf

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	JenkinsUrl string    `yaml:"jenkins-url"`
	DirCI      string    `yaml:"ci-dir"`
	Projects   []Project `yaml:",flow"`
}

type Project struct {
	Path string
	Jobs []Job `yaml:",flow"`
}

type Job struct {
	Branch      string
	Name        string
	Parameters  string
	Token       string
	DiffMatcher string `yaml:"diff-matcher"`
}

func LoadConfig(path *string) (Config, error) {
	fileBytes, readErr := ioutil.ReadFile(*path)
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
