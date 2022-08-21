package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"task/internal/service"
)

type TaskService struct {
	service.UnimplementedTaskServiceServer
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (service *TaskService) TaskCreate(context.Context, *service.TaskRequest) (*service.CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskCreate not implemented")
}
func (service *TaskService) TaskUpdate(context.Context, *service.TaskRequest) (*service.CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskUpdate not implemented")
}
func (service *TaskService) TaskShow(context.Context, *service.TaskRequest) (*service.TaskDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskShow not implemented")
}
func (service *TaskService) TaskDelete(context.Context, *service.TaskRequest) (*service.CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskDelete not implemented")
}
