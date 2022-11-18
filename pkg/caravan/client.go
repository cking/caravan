package caravan

import (
	"errors"
	"fmt"
	"go.uber.org/multierr"

	nomad "github.com/hashicorp/nomad-openapi/clients/go/v1"
	v1 "github.com/hashicorp/nomad-openapi/v1"
)

type Client struct {
	nc *v1.Client
}

func NewClient() (*Client, error) {
	client, err := v1.NewClient()
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}

func (client *Client) GetJob(name string) (*nomad.Job, error) {
	opts := v1.DefaultQueryOpts()
	job, _, err := client.nc.Jobs().GetJob(opts.Ctx(), name)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (client *Client) ListJobs() (map[string]*nomad.Job, error) {
	opts := v1.DefaultQueryOpts()

	joblist, _, err := client.nc.Jobs().GetJobs(opts.Ctx())
	if err != nil {
		return nil, err
	}

	jobs := make(map[string]*nomad.Job)

	err = nil
	for _, job := range *joblist {
		job, e := client.GetJob(*job.Name)
		if multierr.Append(err, e) != nil {
			return nil, err
		}
		jobs[*job.Name] = job
	}

	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (client *Client) ParseJob(job string) (*nomad.Job, error) {
	opts := v1.DefaultQueryOpts()

	parsedJob, err := client.nc.Jobs().Parse(opts.Ctx(), job, false, false)
	if err != nil {
		return nil, err
	}

	return parsedJob, nil
}

// https://github.com/hashicorp/nomad-openapi/
// https://docs.google.com/presentation/d/1h4OOjPFOHbDJsbtuQZRYDjotyBH1YZs7V8L7qmEjRXc/edit#slide=id.gd36c5fdcb4_1_200
func (client *Client) ApplyJob(job *nomad.Job) (string, error) {
	opts := v1.DefaultQueryOpts()

	// Adding metadata to identify the jobs managed by the Nomoporator
	// receive existng metaa
	metadata := job.GetMeta()
	if metadata == nil {
		metadata = make(map[string]string)
	}
	metadata[MetaKeyName] = "true"
	metadata["uid"] = "caravan"
	job.SetMeta(metadata)

	_, _, err := client.nc.Jobs().Plan(opts.Ctx(), job, false)
	if err != nil {
		return "", fmt.Errorf("error while running nomad plan: %w", err)
	}

	res, _, err := client.nc.Jobs().Post(opts.Ctx(), job)
	if err != nil {
		return "", fmt.Errorf("error while running nomad post: %w", err)
	}

	return *res.EvalID, nil
}

func (client *Client) DeleteJob(job *nomad.Job) error {
	opts := v1.DefaultQueryOpts()

	if job.Meta == nil {
		return errors.New("will only delete caravan jobs")
	}

	if _, found := (*job.Meta)[MetaKeyName]; !found {
		return errors.New("will only delete caravan jobs")
	}

	_, _, err := client.nc.Jobs().Delete(opts.Ctx(), job.GetName(), true, true)
	if err != nil {
		return err
	}

	return nil
}
