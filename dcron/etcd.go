package dcron

import (
	"encoding/json"
	etcdc "github.com/coreos/go-etcd/etcd"
)

var machines = []string{"http://127.0.0.1:2379"}
var etcd = NewClient(machines)
var keyspace = "/dcron"

type etcdClient struct {
	Client *etcdc.Client
}

func NewClient(machines []string) *etcdClient {
	return &etcdClient{Client: etcdc.NewClient(machines)}
}

func (e *etcdClient) SetJob(job *Job) error {
	jobJson, _ := json.Marshal(job)
	log.Debugf("Setting etcd key %s: %s", job.Name, string(jobJson))
	if _, err := e.Client.Set(keyspace+"/jobs/"+job.Name, string(jobJson), 0); err != nil {
		return err
	}

	return nil
}

func (e *etcdClient) GetJobs() ([]*Job, error) {
	res, err := e.Client.Get(keyspace+"/jobs/", true, false)
	if err != nil {
		return nil, err
	}

	var jobs []*Job
	for _, node := range res.Node.Nodes {

		log.Debug(*node)
		var job Job
		err := json.Unmarshal([]byte(node.Value), &job)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
		log.Debug(job)
	}
	return jobs, nil
}

func (e *etcdClient) GetJob(name string) (*Job, error) {
	res, err := e.Client.Get(keyspace+"/jobs/"+name, false, false)
	if err != nil {
		return nil, err
	}

	var job Job
	if err = json.Unmarshal([]byte(res.Node.Value), &job); err != nil {
		return nil, err
	}
	log.Debugf("Retrieved job from datastore: %v", job)
	return &job, nil
}

func (e *etcdClient) GetExecutions() ([]*Execution, error) {
	res, err := e.Client.Get(keyspace+"/executions/", true, false)
	if err != nil {
		return nil, err
	}

	var executions []*Execution
	for _, node := range res.Node.Nodes {
		var execution Execution
		err := json.Unmarshal([]byte(node.Value), &execution)
		if err != nil {
			return nil, err
		}
		executions = append(executions, &execution)
	}
	return executions, nil
}

func (e *etcdClient) GetLeader() string {
	res, err := e.Client.Get(keyspace+"/leader", false, false)
	if err != nil {
		log.Debug(err)
		return ""
	}

	log.Debugf("Retrieved leader from datastore: %v", res.Node.Value)
	return res.Node.Value
}
