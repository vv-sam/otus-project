package grpc_services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	api "github.com/vv-sam/otus-project/proto/grpc/pkg"
	"github.com/vv-sam/otus-project/server/internal/model/configuration"
	"github.com/vv-sam/otus-project/server/internal/repository"
	"github.com/vv-sam/otus-project/server/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type configurationRepository interface {
	Get(id uuid.UUID) (*configuration.Factorio, error)
	GetAll() ([]*configuration.Factorio, error)
	Add(configuration *configuration.Factorio) error
	Update(id uuid.UUID, configuration *configuration.Factorio) error
	Delete(id uuid.UUID) error
}

type ConfigurationService struct {
	api.UnimplementedConfigurationServiceServer
	configurationRepository configurationRepository
	validator               *services.Validator
}

func NewConfigurationService(configurationRepository configurationRepository, validator *services.Validator) *ConfigurationService {
	return &ConfigurationService{
		configurationRepository: configurationRepository,
		validator:               validator,
	}
}

func (s *ConfigurationService) GetById(ctx context.Context, req *api.GetConfigByIdRequest) (*api.GetConfigByIdResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	configUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}

	configInfo, err := s.configurationRepository.Get(configUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get configuration: %v", err)
	}

	if configInfo == nil {
		return nil, status.Error(codes.NotFound, "configuration not found")
	}

	protoConfig := convertConfigurationToProto(configInfo)
	return &api.GetConfigByIdResponse{Configuration: protoConfig}, nil
}

func (s *ConfigurationService) GetAll(ctx context.Context, req *api.GetAllConfigurationsRequest) (*api.GetAllConfigurationsResponse, error) {
	configurations, err := s.configurationRepository.GetAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get configurations: %v", err)
	}

	protoConfigurations := make([]*api.FactorioConfig, len(configurations))
	for i, config := range configurations {
		protoConfigurations[i] = convertConfigurationToProto(config)
	}

	return &api.GetAllConfigurationsResponse{Configurations: protoConfigurations}, nil
}

func (s *ConfigurationService) Post(ctx context.Context, req *api.PostConfigurationRequest) (*api.PostConfigurationResponse, error) {
	if req.Configuration == nil {
		return nil, status.Error(codes.InvalidArgument, "configuration is required")
	}

	configInfo, err := convertProtoToConfiguration(req.Configuration)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to convert configuration: %v", err)
	}

	if !s.validator.IsValid(configInfo) {
		return nil, status.Error(codes.InvalidArgument, "invalid configuration")
	}

	if err := s.configurationRepository.Add(configInfo); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add configuration: %v", err)
	}

	return &api.PostConfigurationResponse{}, nil
}

func (s *ConfigurationService) Put(ctx context.Context, req *api.PutConfigurationRequest) (*api.PutConfigurationResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if req.Configuration == nil {
		return nil, status.Error(codes.InvalidArgument, "configuration is required")
	}

	configUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}

	configInfo, err := convertProtoToConfiguration(req.Configuration)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to convert configuration: %v", err)
	}

	// Set the configuration ID from the request
	configInfo.Id = configUUID

	if !s.validator.IsValid(configInfo) {
		return nil, status.Error(codes.InvalidArgument, "invalid configuration")
	}

	if err := s.configurationRepository.Update(configUUID, configInfo); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "configuration not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update configuration: %v", err)
	}

	return &api.PutConfigurationResponse{}, nil
}

func (s *ConfigurationService) Delete(ctx context.Context, req *api.DeleteConfigurationRequest) (*api.DeleteConfigurationResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	configUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}

	if err := s.configurationRepository.Delete(configUUID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "configuration not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete configuration: %v", err)
	}

	return &api.DeleteConfigurationResponse{}, nil
}

func convertConfigurationToProto(config *configuration.Factorio) *api.FactorioConfig {
	return &api.FactorioConfig{
		Base: &api.BaseConfig{
			Id:      config.Id.String(),
			AgentId: config.AgentId.String(),
			Port:    uint32(config.Port),
			Type:    api.ConfigurationType_CONFIGURATION_TYPE_FACTORIO,
		},
	}
}

func convertProtoToConfiguration(protoConfig *api.FactorioConfig) (*configuration.Factorio, error) {
	if protoConfig.Base == nil {
		return nil, errors.New("base configuration is required")
	}

	configUUID, err := uuid.Parse(protoConfig.Base.Id)
	if err != nil {
		return nil, err
	}

	agentUUID, err := uuid.Parse(protoConfig.Base.AgentId)
	if err != nil {
		return nil, err
	}

	config := &configuration.Factorio{}
	config.Id = configUUID
	config.AgentId = agentUUID
	config.Port = uint16(protoConfig.Base.Port)
	config.Type = configuration.CONFIGURATION_TYPE_FACTORIO

	return config, nil
}
