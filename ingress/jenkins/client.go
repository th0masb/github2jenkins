package jenkins

import (
	"fmt"
	"net/http"

	"github.com/th0masb/github2jenkins/tools"
)

const (
	job             string = "job"
	build           string = "build"
	buildWithParams string = "buildWithParameters"
	token           string = "token"
	cause           string = "cause"
)

// Client Execute trigger request to Jenkins instance
type Client struct {
	config   ClientConfig
	delegate tools.HTTPRequester
}

// ClientConfig Configure the client
type ClientConfig struct {
	JenkinsURL string
}

// TriggerJobRequest Request information for choosing right job
type TriggerJobRequest struct {
	JobName    string
	Cause      string
	Token      string
	Parameters map[string]string
}

// TriggerJob Instruct the client to start the job described by the request
func (c *Client) TriggerJob(req TriggerJobRequest) error {
	url := c.buildRequestURL(&req)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := c.delegate.Do(httpReq)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if tools.IsErrorCode(resp.StatusCode) {
		return fmt.Errorf(
			"Recieved error response status: %s for request: %+v",
			resp.Status,
			req,
		)
	}
	return nil
}

func (c *Client) buildRequestURL(req *TriggerJobRequest) string {
	var buildArg string
	if len(req.Parameters) > 0 {
		buildArg = buildWithParams
	} else {
		buildArg = build
	}

	params := fmt.Sprintf("%s=%s&%s=%s", token, req.Token, cause, req.Cause)
	for k, v := range req.Parameters {
		params += fmt.Sprintf("&%s=%s", k, v)
	}

	return fmt.Sprintf(
		"%s/%s/%s/%s?%s",
		c.config.JenkinsURL,
		job,
		req.JobName,
		buildArg,
		params,
	)
}
