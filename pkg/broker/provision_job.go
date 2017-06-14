package broker

import (
	"encoding/json"

	logging "github.com/op/go-logging"
	"github.com/openshift/ansible-service-broker/pkg/apb"
	"github.com/pborman/uuid"
)

type ProvisionJob struct {
	instanceuuid  uuid.UUID
	spec          *apb.Spec
	context       *apb.Context
	parameters    *apb.Parameters
	clusterConfig apb.ClusterConfig
	log           *logging.Logger
}

type ProvisionMsg struct {
	InstanceUUID string `json:"instance_uuid"`
	JobToken     string `json:"job_token"`
	SpecId       string `json:"spec_id"`
	PodName      string `json:"podname"`
	Msg          string `json:"msg"`
	Error        string `json:"error"`
}

func (m ProvisionMsg) Render() string {
	render, _ := json.Marshal(m)
	return string(render)
}

func NewProvisionJob(
	instanceuuid uuid.UUID, spec *apb.Spec, context *apb.Context,
	parameters *apb.Parameters, clusterConfig apb.ClusterConfig,
	log *logging.Logger,
) *ProvisionJob {
	return &ProvisionJob{
		instanceuuid:  instanceuuid,
		spec:          spec,
		context:       context,
		parameters:    parameters,
		clusterConfig: clusterConfig,
		log:           log}
}

func (p *ProvisionJob) Run(token string, msgBuffer chan<- WorkMsg) {
	output, err := apb.Provision(p.spec, p.context, p.parameters, p.clusterConfig, p.log)
	if err != nil {
		p.log.Error("broker::Provision error occurred.")
		p.log.Error("%s", err.Error())
		// send error message
		// can't have an error type in a struct you want marshalled
		// https://github.com/golang/go/issues/5161
		msgBuffer <- ProvisionMsg{InstanceUUID: p.instanceuuid.String(),
			JobToken: token, SpecId: p.spec.Id, PodName: "", Msg: "", Error: err.Error()}
		return
	}

	// save off podname
	podname, _ := apb.GetPodName(output, p.log)
	msgBuffer <- ProvisionMsg{InstanceUUID: p.instanceuuid.String(),
		JobToken: token, SpecId: p.spec.Id, PodName: podname, Msg: "", Error: ""}

	// need to get the pod name for the job state
	extCreds, extErr := apb.ExtractCredentials(podname, p.log)
	if extErr != nil {
		p.log.Error("broker::Provision extError occurred.")
		p.log.Error("%s", extErr.Error())
		// send extError message
		// can't have an extError type in a struct you want marshalled
		// https://github.com/golang/go/issues/5161
		msgBuffer <- ProvisionMsg{InstanceUUID: p.instanceuuid.String(),
			JobToken: token, SpecId: p.spec.Id, PodName: podname, Msg: "", Error: extErr.Error()}
		return
	}

	// send creds
	jsonmsg, _ := json.Marshal(extCreds)
	p.log.Debug("sending message to channel")
	msgBuffer <- ProvisionMsg{InstanceUUID: p.instanceuuid.String(),
		JobToken: token, SpecId: p.spec.Id, PodName: podname, Msg: string(jsonmsg), Error: ""}
}
