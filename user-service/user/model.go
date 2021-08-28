package user

import (
	"context"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Image     string `json:"image"`
	Role      string `json:"role"`
	ResetCode string `json:"resetCode"`
}

type UserProfile struct {
	Email   string `json:"email" gorm:"unique"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Image   string `json:"image"`
	Role    string `json:"role"`
}

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Repository interface {
	CreateUser(ctx context.Context, user User) error
	LoginUser(ctx context.Context, credentials LoginInfo) (string, error)
	GetUser(ctx context.Context, email string) (UserProfile, error)
	ChangePassword(ctx context.Context, email string, resetCode string, newPassword string) error
	EditProfile(ctx context.Context, email string, name string, surname string, image string) error
	ChangeRole(ctx context.Context, email string) error
}
