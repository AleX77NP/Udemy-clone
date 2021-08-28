package user

import (
	"context"
	"errors"
)

type Service interface {
	CreateUser(ctx context.Context, email string, password string, name string, surname string, image string, role string) (string, error)
	LoginUser(ctx context.Context, email string, password string) (string, error)
	GetUser(ctx context.Context, email string) (UserProfile, error)
	ChangePassword(ctx context.Context, email string, resetCode string, newPassword string) (string, error)
	EditProfile(ctx context.Context, email string, name string, surname string, image string) (string, error)
	ChangeRole(ctx context.Context, email string) (string, error)
}

var (
	NotFound      = errors.New("User not found.")
	WrongPassword = errors.New("Incorrect password")
	UniqueEmail   = errors.New("Email is already taken")
	InvalidData   = errors.New("Invalid data provided. Cannot parse user.")
	WrongCode     = errors.New("Reset code is not valid.")
)
