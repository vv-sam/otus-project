package grpc_services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	api "github.com/vv-sam/otus-project/proto/grpc/pkg"
	"github.com/vv-sam/otus-project/server/internal/model/agent"
	"github.com/vv-sam/otus-project/server/internal/model/metrics"
	"github.com/vv-sam/otus-project/server/internal/model/task"
	"github.com/vv-sam/otus-project/server/internal/repository"
	"github.com/vv-sam/otus-project/server/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type agentRepository interface {
	Get(id uuid.UUID) (*agent.Info, error)
	GetAll() ([]*agent.Info, error)
	Add(agent *agent.Info) error
	Update(id uuid.UUID, agent *agent.Info) error
	Delete(id uuid.UUID) error
}

type AgentService struct {
	api.UnimplementedAgentServiceServer
	agentRepository agentRepository
	validator       *services.Validator
}

func NewAgentService(agentRepository agentRepository, validator *services.Validator) *AgentService {
	return &AgentService{
		agentRepository: agentRepository,
		validator:       validator,
	}
}

func (s *AgentService) GetById(ctx context.Context, req *api.GetAgentByIdRequest) (*api.GetAgentByIdResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	agentUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}

	agentInfo, err := s.agentRepository.Get(agentUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get agent: %v", err)
	}

	if agentInfo == nil {
		return nil, status.Error(codes.NotFound, "agent not found")
	}

	protoAgent := convertAgentToProto(agentInfo)
	return &api.GetAgentByIdResponse{Agent: protoAgent}, nil
}

func (s *AgentService) GetAll(ctx context.Context, req *api.GetAllAgentsRequest) (*api.GetAllAgentsResponse, error) {
	agents, err := s.agentRepository.GetAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get agents: %v", err)
	}

	protoAgents := make([]*api.AgentInfo, len(agents))
	for i, agent := range agents {
		protoAgents[i] = convertAgentToProto(agent)
	}

	return &api.GetAllAgentsResponse{Agents: protoAgents}, nil
}

func (s *AgentService) Post(ctx context.Context, req *api.PostAgentRequest) (*api.PostAgentResponse, error) {
	if req.Agent == nil {
		return nil, status.Error(codes.InvalidArgument, "agent is required")
	}

	agentInfo, err := convertProtoToAgent(req.Agent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to convert agent: %v", err)
	}

	if !s.validator.IsValid(agentInfo) {
		return nil, status.Error(codes.InvalidArgument, "invalid agent")
	}

	if err := s.agentRepository.Add(agentInfo); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add agent: %v", err)
	}

	return &api.PostAgentResponse{}, nil
}

func (s *AgentService) Put(ctx context.Context, req *api.PutAgentRequest) (*api.PutAgentResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if req.Agent == nil {
		return nil, status.Error(codes.InvalidArgument, "agent is required")
	}

	agentUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}

	agentInfo, err := convertProtoToAgent(req.Agent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to convert agent: %v", err)
	}

	// Set the agent ID from the request
	agentInfo.AgentId = agentUUID

	if !s.validator.IsValid(agentInfo) {
		return nil, status.Error(codes.InvalidArgument, "invalid agent")
	}

	if err := s.agentRepository.Update(agentUUID, agentInfo); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "agent not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update agent: %v", err)
	}

	return &api.PutAgentResponse{}, nil
}

func (s *AgentService) Delete(ctx context.Context, req *api.DeleteAgentRequest) (*api.DeleteAgentResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	agentUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}

	if err := s.agentRepository.Delete(agentUUID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "agent not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete agent: %v", err)
	}

	return &api.DeleteAgentResponse{}, nil
}

// Helper functions for conversion between internal models and protobuf messages

func convertAgentToProto(agentInfo *agent.Info) *api.AgentInfo {
	protoTasks := make([]*api.Task, len(agentInfo.CurrentTasks))
	for i, task := range agentInfo.CurrentTasks {
		protoTasks[i] = convertTaskToProto(&task)
	}

	return &api.AgentInfo{
		AgentId:      agentInfo.AgentId.String(),
		Status:       int32(agentInfo.Status),
		CurrentTasks: protoTasks,
		Metrics:      convertMetricsToProto(&agentInfo.Metrics),
	}
}

func convertProtoToAgent(protoAgent *api.AgentInfo) (*agent.Info, error) {
	agentUUID, err := uuid.Parse(protoAgent.AgentId)
	if err != nil {
		return nil, err
	}

	tasks := make([]task.Task, len(protoAgent.CurrentTasks))
	for i, protoTask := range protoAgent.CurrentTasks {
		convertedTask, err := convertProtoToTask(protoTask)
		if err != nil {
			return nil, err
		}
		tasks[i] = *convertedTask
	}

	var hostMetrics metrics.HostMetrics
	if protoAgent.Metrics != nil {
		hostMetrics = *convertProtoToMetrics(protoAgent.Metrics)
	}

	return &agent.Info{
		AgentId:      agentUUID,
		Status:       int16(protoAgent.Status),
		CurrentTasks: tasks,
		Metrics:      hostMetrics,
	}, nil
}

func convertTaskToProto(taskInfo *task.Task) *api.Task {
	return &api.Task{
		Id:     taskInfo.Id.String(),
		Status: convertTaskStatusToProto(taskInfo.Status),
		Type:   taskInfo.Type,
	}
}

func convertProtoToTask(protoTask *api.Task) (*task.Task, error) {
	taskUUID, err := uuid.Parse(protoTask.Id)
	if err != nil {
		return nil, err
	}

	return &task.Task{
		Id:     taskUUID,
		Status: convertProtoToTaskStatus(protoTask.Status),
		Type:   protoTask.Type,
	}, nil
}

func convertTaskStatusToProto(status int16) api.TaskStatus {
	switch status {
	case task.STATUS_QUEUED:
		return api.TaskStatus_TASK_STATUS_QUEUED
	case task.STATUS_IN_PROGRESS:
		return api.TaskStatus_TASK_STATUS_IN_PROGRESS
	case task.STATUS_OK:
		return api.TaskStatus_TASK_STATUS_OK
	case task.STATUS_DELETED:
		return api.TaskStatus_TASK_STATUS_DELETED
	default:
		return api.TaskStatus_TASK_STATUS_UNSPECIFIED
	}
}

func convertProtoToTaskStatus(status api.TaskStatus) int16 {
	switch status {
	case api.TaskStatus_TASK_STATUS_QUEUED:
		return task.STATUS_QUEUED
	case api.TaskStatus_TASK_STATUS_IN_PROGRESS:
		return task.STATUS_IN_PROGRESS
	case api.TaskStatus_TASK_STATUS_OK:
		return task.STATUS_OK
	case api.TaskStatus_TASK_STATUS_DELETED:
		return task.STATUS_DELETED
	default:
		return task.STATUS_QUEUED
	}
}

func convertMetricsToProto(hostMetrics *metrics.HostMetrics) *api.HostMetrics {
	return &api.HostMetrics{
		CpuUsage:     hostMetrics.CpuUsage,
		RamAvailable: hostMetrics.RamAvailable,
		RamTotal:     hostMetrics.RamTotal,
	}
}

func convertProtoToMetrics(protoMetrics *api.HostMetrics) *metrics.HostMetrics {
	return &metrics.HostMetrics{
		CpuUsage:     protoMetrics.CpuUsage,
		RamAvailable: protoMetrics.RamAvailable,
		RamTotal:     protoMetrics.RamTotal,
	}
}
