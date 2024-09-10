package service

import (
	"config-server2/internal/common"
	"config-server2/internal/entity"
	"config-server2/internal/protov2"
	"config-server2/internal/repository"
)

func CreateAgentGroup(req *protov2.CreateAgentGroupRequest, res *protov2.CreateAgentGroupResponse) error {
	requestId := req.RequestId
	if requestId == nil {
		return common.ValidateErrorWithMsg("required fields requestId could not be null")
	}

	agentGroup := req.AgentGroup
	if agentGroup == nil {
		return common.ValidateErrorWithMsg("required field agentGroup could not be null")
	}

	res.RequestId = requestId
	group := entity.ParseProtoAgentGroupTag2AgentGroup(agentGroup)
	err := repository.CreateAgentGroup(group)
	return common.SystemError(err)
}

func UpdateAgentGroup(req *protov2.UpdateAgentGroupRequest, res *protov2.UpdateAgentGroupResponse) error {
	requestId := req.RequestId
	if requestId == nil {
		return common.ValidateErrorWithMsg("required fields requestId could not be null")
	}

	agentGroup := req.AgentGroup
	if agentGroup == nil {
		return common.ValidateErrorWithMsg("required field agentGroup could not be null")
	}

	res.RequestId = requestId
	group := entity.ParseProtoAgentGroupTag2AgentGroup(agentGroup)
	err := repository.UpdateAgentGroup(group)
	return common.SystemError(err)
}

func DeleteAgentGroup(req *protov2.DeleteAgentGroupRequest, res *protov2.DeleteAgentGroupResponse) error {
	requestId := req.RequestId
	if requestId == nil {
		return common.ValidateErrorWithMsg("required fields requestId could not be null")
	}

	agentGroupName := req.GroupName
	if agentGroupName == "" {
		return common.ValidateErrorWithMsg("required field groupName could not be null")
	}
	res.RequestId = requestId
	if req.GroupName == entity.AgentGroupDefaultValue {
		return common.ServerErrorWithMsg("%s can not be deleted", entity.AgentGroupDefaultValue)
	}
	err := repository.DeleteAgentGroup(agentGroupName)
	return common.SystemError(err)
}

func GetAgentGroup(req *protov2.GetAgentGroupRequest, res *protov2.GetAgentGroupResponse) error {
	requestId := req.RequestId
	if requestId == nil {
		return common.ValidateErrorWithMsg("required fields requestId could not be null")
	}

	agentGroupName := req.GroupName
	if agentGroupName == "" {
		return common.ValidateErrorWithMsg("required field groupName could not be null")
	}

	res.RequestId = requestId
	agentGroup, err := repository.GetAgentGroupDetail(agentGroupName, false, false)
	if err != nil {
		return common.SystemError(err)
	}
	res.AgentGroup = agentGroup.Parse2ProtoAgentGroupTag()
	return nil
}

func ListAgentGroups(req *protov2.ListAgentGroupsRequest, res *protov2.ListAgentGroupsResponse) error {
	requestId := req.RequestId
	if requestId == nil {
		return common.ValidateErrorWithMsg("required fields requestId could not be null")
	}

	res.RequestId = requestId
	agentGroups, err := repository.GetAllAgentGroup()
	if err != nil {
		return common.SystemError(err)
	}
	protoAgentGroups := make([]*protov2.AgentGroupTag, 0)
	for _, agentGroup := range agentGroups {
		protoAgentGroups = append(protoAgentGroups, agentGroup.Parse2ProtoAgentGroupTag())
	}
	res.AgentGroups = protoAgentGroups
	return nil
}

func GetAppliedAgentGroupsForPipelineConfigName(req *protov2.GetAppliedAgentGroupsRequest, res *protov2.GetAppliedAgentGroupsResponse) error {
	requestId := req.RequestId
	if requestId == nil {
		return common.ValidateErrorWithMsg("required fields requestId could not be null")
	}
	configName := req.ConfigName
	if configName == "" {
		return common.ValidateErrorWithMsg("required fields configName could not be null")
	}

	res.RequestId = requestId

	groupNames, err := repository.GetAppliedAgentGroupForPipelineConfigName(configName)
	if err != nil {
		return common.SystemError(err)
	}
	res.AgentGroupNames = groupNames
	return nil
}

func GetAppliedAgentGroupsForInstanceConfigName(req *protov2.GetAppliedAgentGroupsRequest, res *protov2.GetAppliedAgentGroupsResponse) error {
	requestId := req.RequestId
	if requestId == nil {
		return common.ValidateErrorWithMsg("required fields requestId could not be null")
	}
	configName := req.ConfigName
	if configName == "" {
		return common.ValidateErrorWithMsg("required fields configName could not be null")
	}

	res.RequestId = requestId

	groupNames, err := repository.GetAppliedAgentGroupForInstanceConfigName(configName)
	if err != nil {
		return common.SystemError(err)
	}
	res.AgentGroupNames = groupNames
	return nil
}
