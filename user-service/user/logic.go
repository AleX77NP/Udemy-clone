package user

import (
	"context"
	"crypto/rand"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type service struct {
	repository Repository
	logger     log.Logger
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func NewService(rep Repository, logger log.Logger) Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

func (s service) CreateUser(ctx context.Context, email string, password string, name string, surname string, image string, role string) (string, error) {
	logger := log.With(s.logger, "method", "CreateUser")

	var code = EncodeToString(4)

	user := User{
		Email:     email,
		Password:  password,
		Name:      name,
		Surname:   surname,
		Image:     image,
		Role:      role,
		ResetCode: code,
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("Create user", user.ID)

	return code, nil
}

func (s service) LoginUser(ctx context.Context, email string, password string) (string, error) {
	logger := log.With(s.logger, "method", "LoginUser")

	var jwt string

	credentials := LoginInfo{
		Email:    email,
		Password: password,
	}

	if token, err := s.repository.LoginUser(ctx, credentials); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	} else {
		jwt = token
	}

	logger.Log("Login user", email)

	return jwt, nil
}

func (s service) GetUser(ctx context.Context, email string) (UserProfile, error) {
	logger := log.With(s.logger, "method", "GetUser")

	if userProfile, err := s.repository.GetUser(ctx, email); err != nil {
		return UserProfile{}, err
	} else {
		logger.Log("Get user", email)
		return userProfile, nil
	}
}

func (s service) ChangePassword(ctx context.Context, email string, resetCode string, newPassword string) (string, error) {
	logger := log.With(s.logger, "method", "ChangePassword")
	if err := s.repository.ChangePassword(ctx, email, resetCode, newPassword); err != nil {
		return "", err
	}

	logger.Log("Change password")
	return "Password reset success.", nil
}

func (s service) EditProfile(ctx context.Context, email string, name string, surname string, image string) (string, error) {
	logger := log.With(s.logger, "method", "EditProfile")

	if err := s.repository.EditProfile(ctx, email, name, surname, image); err != nil {
		return "", err
	}

	logger.Log("Edit profile", email)
	return "Profile edited.", nil
}

func (s service) ChangeRole(ctx context.Context, email string) (string, error) {
	logger := log.With(s.logger, "method", "ChangeRole")

	if err := s.repository.ChangeRole(ctx, email); err != nil {
		return "", err
	}

	logger.Log("Role Changed", email)
	return "Role change success", nil
}
