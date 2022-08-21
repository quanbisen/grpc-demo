package handler

import (
	"context"
	"user/common/res"
	"user/internal/repository"
	"user/internal/service"
)

type UserService struct {
	service.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) UserLogin(ctx context.Context, req *service.UserRequest) (*service.UserDetailResponse, error) {
	var user repository.User
	resp := new(service.UserDetailResponse)
	resp.Code = res.Success
	err := user.ShowUserInfo(req)
	if err != nil {
		resp.Code = res.Error
		return resp, err
	}
	resp.UserDetail = repository.BuildUser(user)
	return resp, nil
}
func (s *UserService) UserRegister(ctx context.Context, req *service.UserRequest) (*service.UserDetailResponse, error) {
	var user repository.User
	resp := new(service.UserDetailResponse)
	resp.Code = res.Success
	err := user.UserCreate(req)
	if err != nil {
		resp.Code = res.Error
		return resp, err
	}
	resp.UserDetail = repository.BuildUser(user)
	return resp, nil
}
