package capability

import (
	"config-server2/internal/common"
	proto "config-server2/internal/common/protov2"
	"config-server2/internal/config"
	"config-server2/internal/entity"
	"config-server2/internal/manager"
)

type ServerAction struct {
	Base
	run func(*proto.HeartbeatRequest, *proto.HeartbeatResponse) error
}

func (a ServerAction) Action(req *proto.HeartbeatRequest, res *proto.HeartbeatResponse) error {
	serverCapabilityType := a.Value
	hasServerCapability, ok := config.ServerConfigInstance.Capabilities[serverCapabilityType]
	if !ok {
		return common.ServerErrorWithMsg("error serverCapability type(%s)", serverCapabilityType)
	}

	if hasServerCapability {
		err := a.run(req, res)
		return err
	}
	return nil
}

func (a ServerAction) UpdateCapabilities(res *proto.HeartbeatResponse) {
	res.Capabilities = res.Capabilities | uint64(a.Code)
}

var (
	UnspecifiedServerAction = &ServerAction{
		Base: ServerUnspecified,
		run:  UnspecifiedServerCapabilityRun,
	}

	RememberAttributeAction = &ServerAction{
		Base: RememberAttribute,
		run:  RememberAttributeCapabilityRun,
	}
	RememberPipelineConfigStatusAction = &ServerAction{
		Base: RememberPipelineConfigStatus,
		run:  RememberPipelineConfigStatusCapabilityRun,
	}
	RememberProcessConfigStatusAction = &ServerAction{
		Base: RememberProcessConfigStatus,
		run:  RememberProcessConfigStatusCapabilityRun,
	}
	RememberCustomCommandStatusAction = &ServerAction{
		Base: RememberCustomCommandStatus,
		run:  RememberCustomCommandStatusRun,
	}
)

var ServerActionList = []*ServerAction{
	UnspecifiedServerAction,
	RememberAttributeAction,
	RememberPipelineConfigStatusAction,
	RememberProcessConfigStatusAction,
	RememberCustomCommandStatusAction,
}

func UnspecifiedServerCapabilityRun(req *proto.HeartbeatRequest, res *proto.HeartbeatResponse) error {
	return nil
}

func RememberAttributeCapabilityRun(req *proto.HeartbeatRequest, res *proto.HeartbeatResponse) error {
	attributes := req.Attributes
	if attributes == nil {
		return nil
	}
	agent := &entity.Agent{}
	agent.Attributes = entity.ProtoAgentAttributesParse2AgentAttributes(attributes)
	err := s.UpdateAgentById(agent)
	return err
}

func RememberPipelineConfigStatusCapabilityRun(req *proto.HeartbeatRequest, res *proto.HeartbeatResponse) error {
	configs := req.PipelineConfigs
	if configs == nil {
		return nil
	}
	agentPipelineConfigs := entity.ProtoConfigInfoParse2AgentPipelineConfig(string(req.InstanceId), req.PipelineConfigs)
	err := manager.CreateOrUpdateAgentPipelineConfig(agentPipelineConfigs)
	return err
}

func RememberProcessConfigStatusCapabilityRun(req *proto.HeartbeatRequest, res *proto.HeartbeatResponse) error {
	configs := req.ProcessConfigs
	if configs == nil {
		return nil
	}
	agentProcessConfigs := entity.ProtoConfigInfoParse2AgentProcessConfig(string(req.InstanceId), req.ProcessConfigs)
	err := manager.CreateOrUpdateAgentProcessConfig(agentProcessConfigs)
	return err
}

func RememberCustomCommandStatusRun(req *proto.HeartbeatRequest, res *proto.HeartbeatResponse) error {
	return nil
}