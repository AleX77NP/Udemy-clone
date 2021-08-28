package user

import (
	"github.com/go-kit/kit/endpoint"

	"context"
)

type Endpoints struct {
	CreateUser     endpoint.Endpoint
	LoginUser      endpoint.Endpoint
	GetUser        endpoint.Endpoint
	ChangePassword endpoint.Endpoint
	EditProfile    endpoint.Endpoint
	ChangeRole     endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateUser:     makeCreateUserEndoint(s),
		LoginUser:      makeLoginUserEndoint(s),
		GetUser:        makeGetUserEndpoint(s),
		ChangePassword: makeChangePasswordEndpoint(s),
		EditProfile:    makeEditProfileEndpoint(s),
		ChangeRole:     makeChangeRoleEndpoint(s),
	}
}

func makeCreateUserEndoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		rc, err := s.CreateUser(ctx, req.Email, req.Password, req.Name, req.Surname, req.Image, req.Role)
		return CreateUserResponse{ResetCode: rc}, err
	}
}

func makeLoginUserEndoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateLoginUserRequest)
		token, err := s.LoginUser(ctx, req.Email, req.Password)
		return CreateLoginUserResponse{Token: token}, err
	}
}

func makeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateGetUserRequest)
		userProfile, err := s.GetUser(ctx, req.Email)
		return CreateGetUserResponse{UserProfile: userProfile}, err
	}
}

func makeChangePasswordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateChangePasswordRequest)
		ok, err := s.ChangePassword(ctx, req.Email, req.ResetCode, req.NewPassword)
		return CreateChangePasswordResponse{Ok: ok}, err
	}
}

func makeEditProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateEditProfileRequest)
		ok, err := s.EditProfile(ctx, req.Email, req.Name, req.Surname, req.Image)
		return CreateEditProfileResponse{Ok: ok}, err
	}
}

func makeChangeRoleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateChangeRoleRequest)
		ok, err := s.ChangeRole(ctx, req.Email)
		return CreateChangeRoleResponse{Ok: ok}, err
	}
}
