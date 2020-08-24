package jenkins

import (
	"fmt"

	"github.com/th0masb/github2jenkins/g2j"
)

const (
	allBranchMatcher  rune = '*'
	negationMatchRune rune = '^'
)

type projectAndChangedFiles struct {
	project      g2j.Project
	changedFiles []string
}

type jobWithMatcher struct {
	job g2j.Job
}

// TriggerService Decides which jobs should be triggered from a github event
type TriggerService struct {
	config TriggerServiceConfig
	client *Client
}

// TriggerServiceConfig Configures the trigger service
type TriggerServiceConfig struct {
	Repositories []g2j.Repository
}

// TriggerJobs Trigger the relevant jobs for the given repo and the files changed
func (s TriggerService) TriggerJobs(
	repoName string,
	ref string,
	changedFiles []string,
) error {
	//	repo, err := s.findRepository(repoName)
	//	if err != nil {
	//		return err
	//	}
	//
	//	for _, p := range repo.Projects {
	//		applicableFiles := findApplicableFiles(p, changedFiles)
	//	}
	//
	//	// Assign each project a subslice of the changedFiles according to whether
	//	// the file path starts with the project path
	//
	//	// Remove the project path prefix from each of the changed files in the
	//	// assigned slice
	//
	//	// For each job under each (project, changedFiles) pair check if it matches
	//	// the ref and the set of changed files, if it does then trigger the job

	return nil
}

//func findApplicableFiles(project *g2j.Project, changedFiles []string) []string {
//	dest := []string{}
//	for _, cf := range changedFiles {
//		if strings.HasPrefix(cf, project.Path) {
//
//		}
//
//	}
//
//}

func (s TriggerService) findRepository(repoName string) (*g2j.Repository, error) {
	for i := range s.config.Repositories {
		repo := &s.config.Repositories[i]
		if repo.Name == repoName {
			return repo, nil
		}
	}
	return nil, fmt.Errorf("Repository %s is not registered", repoName)
}
