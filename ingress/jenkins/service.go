package jenkins

// Service Decides which jobs should be triggered from a github event
type Service struct {
	// need client
}

// TriggerJobs Trigger the relevant jobs for the given repo and the files changed
func (s Service) TriggerJobs(repo string, changedFiles []string) error {
	return nil
}
