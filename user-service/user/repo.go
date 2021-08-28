package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var RepoErr = errors.New("Unable to handle Your request.")

var (
	DBConn *gorm.DB
)

type repo struct {
	db     *gorm.DB
	logger log.Logger
}

func NewRepo(db *gorm.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (repo *repo) CreateUser(ctx context.Context, user User) error {
	db := repo.db

	if user.Email == "" || user.Password == "" || user.Name == "" || user.Surname == "" || user.Role == "" {
		return RepoErr
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return RepoErr
	}

	user.Password = string(hashed)

	if err := db.Create(&user).Error; err != nil {
		return UniqueEmail
	}

	return nil
}

func (repo *repo) LoginUser(ctx context.Context, credentials LoginInfo) (string, error) {
	db := repo.db

	if credentials.Email == "" || credentials.Password == "" {
		return "", RepoErr
	}

	user := new(User)
	if err := db.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		return "", NotFound
	}

	fmt.Println(user.Role)

	if encryptionErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); encryptionErr != nil {
		return "", WrongPassword
	}

	if token, err := generateToken(credentials.Email, user.Role); err != nil {
		return "", errors.New("Server error.")
	} else {
		return token, nil
	}

}

func (repo *repo) GetUser(ctx context.Context, email string) (UserProfile, error) {
	db := repo.db
	var user User

	if email == "" {
		return UserProfile{}, InvalidData
	}

	if result := db.Where("email = ?", email).First(&user); result.Error != nil {
		return UserProfile{}, NotFound
	}
	userProfile := UserProfile{
		Email:   email,
		Name:    user.Name,
		Surname: user.Surname,
		Image:   user.Image,
		Role:    user.Role,
	}
	return userProfile, nil
}

func (repo *repo) ChangePassword(ctx context.Context, email string, resetCode string, newPassword string) error {
	db := repo.db
	var user User

	if email == "" || resetCode == "" || newPassword == "" {
		return InvalidData
	}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return NotFound
	}
	if user.ResetCode != resetCode {
		return WrongCode
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return RepoErr
	}

	user.Password = string(hashed)
	if err := db.Save(&user).Error; err != nil {
		return RepoErr
	}
	return nil
}

func (repo *repo) EditProfile(ctx context.Context, email string, name string, surname string, image string) error {
	db := repo.db
	var user User

	if email == "" || name == "" || surname == "" || image == "" {
		return InvalidData
	}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return NotFound
	}

	if err := db.Model(&user).Updates(User{Name: name, Surname: surname, Image: image}).Error; err != nil {
		return RepoErr
	}
	return nil
}

func (repo *repo) ChangeRole(ctx context.Context, email string) error {
	db := repo.db
	var user User

	if email == "" {
		return InvalidData
	}

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return NotFound
	}

	user.Role = "teacher"

	if err := db.Save(&user).Error; err != nil {
		return RepoErr
	}

	return nil
}
