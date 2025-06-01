package grpc_services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	api "github.com/vv-sam/otus-project/proto/grpc/pkg"
	"github.com/vv-sam/otus-project/server/internal/model/task"
	"github.com/vv-sam/otus-project/server/internal/repository"
	"github.com/vv-sam/otus-project/server/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type tasksRepository interface {
	Get(id uuid.UUID) (*task.Task, error)
	GetAll() ([]*task.Task, error)
	Add(task *task.Task) error
	Update(id uuid.UUID, task *task.Task) error
	Delete(id uuid.UUID) error
}

type TaskService struct {
	api.UnimplementedTaskServiceServer
	tasksRepository tasksRepository
	validator       *services.Validator
}

func NewTaskService(tasksRepository tasksRepository, validator *services.Validator) *TaskService {
	return &TaskService{
		tasksRepository: tasksRepository,
		validator:       validator,
	}
}

func (s *TaskService) GetById(ctx context.Context, req *api.GetTaskByIdRequest) (*api.GetTaskByIdResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	taskUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}

	taskInfo, err := s.tasksRepository.Get(taskUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get task: %v", err)
	}

	if taskInfo == nil {
		return nil, status.Error(codes.NotFound, "task not found")
	}

	protoTask := convertTaskToProto(taskInfo)
	return &api.GetTaskByIdResponse{Task: protoTask}, nil
}

func (s *TaskService) GetAll(ctx context.Context, req *api.GetAllTasksRequest) (*api.GetAllTasksResponse, error) {
	tasks, err := s.tasksRepository.GetAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tasks: %v", err)
	}

	protoTasks := make([]*api.Task, len(tasks))
	for i, task := range tasks {
		protoTasks[i] = convertTaskToProto(task)
	}

	return &api.GetAllTasksResponse{Tasks: protoTasks}, nil
}

func (s *TaskService) Post(ctx context.Context, req *api.PostTaskRequest) (*api.PostTaskResponse, error) {
	if req.Task == nil {
		return nil, status.Error(codes.InvalidArgument, "task is required")
	}

	taskInfo, err := convertProtoToTask(req.Task)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to convert task: %v", err)
	}

	if !s.validator.IsValid(taskInfo) {
		return nil, status.Error(codes.InvalidArgument, "invalid task")
	}

	if err := s.tasksRepository.Add(taskInfo); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add task: %v", err)
	}

	return &api.PostTaskResponse{}, nil
}

func (s *TaskService) Put(ctx context.Context, req *api.PutTaskRequest) (*api.PutTaskResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if req.Task == nil {
		return nil, status.Error(codes.InvalidArgument, "task is required")
	}

	taskUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}

	taskInfo, err := convertProtoToTask(req.Task)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to convert task: %v", err)
	}

	// Set the task ID from the request
	taskInfo.Id = taskUUID

	if !s.validator.IsValid(taskInfo) {
		return nil, status.Error(codes.InvalidArgument, "invalid task")
	}

	if err := s.tasksRepository.Update(taskUUID, taskInfo); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update task: %v", err)
	}

	return &api.PutTaskResponse{}, nil
}

func (s *TaskService) Delete(ctx context.Context, req *api.DeleteTaskRequest) (*api.DeleteTaskResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	taskUUID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}

	if err := s.tasksRepository.Delete(taskUUID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete task: %v", err)
	}

	return &api.DeleteTaskResponse{}, nil
}
