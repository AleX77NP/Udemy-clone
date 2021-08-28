package course

import (
	"context"
	"errors"
)

type Service interface {
	GetCourses(ctx context.Context, page int) (CoursesRes, error)
	GetCoursesByCategory(ctx context.Context, category string) ([]Course, error)
	GetCourseById(ctx context.Context, id uint) (CourseWithRating, error)
	CreateCourse(ctx context.Context, title string, description string, price float32, image string, link string, author string, authorRef string, category string) (string, error)
	UpdateCourse(ctx context.Context, id uint, title string, description string, price float32, image string, link string, authorRef string) (string, error)
	DeleteCourse(ctx context.Context, id uint, authorRef string) (string, error)
	CreateOrder(ctx context.Context, courses []Course, user string, totalPrice float32) (string, error)
	RateCourse(ctx context.Context, id uint, rating int32, user string) (string, error)
	GetUserOrders(ctx context.Context, email string) ([]Order, error)
	GetTeacherCourses(ctx context.Context, email string) ([]Course, error)
}

var (
	NotFound = errors.New("Course not found.")
	InvalidData = errors.New("Invalid data")
	Unauthorized = errors.New("Unauthorized request.")
	GetRatingError = errors.New("Couldn't get rating.")
	AddRatingError = errors.New("Couldn't add rating.")
)