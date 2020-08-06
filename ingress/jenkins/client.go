package jenkins

type Client struct {
}

type TriggerJobRequest struct {
	JobName    string
	Cause      string
	Token      string
	Parameters map[string]string
}

func (c Client) TriggerJob(req TriggerJobRequest) error {
	return nil
}
