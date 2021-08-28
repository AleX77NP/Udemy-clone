package course

import (
	"gorm.io/gorm"
	"context"
)

type Course struct {
	gorm.Model
	Title string
	Description string
	Price float32
	Image string
	Link string
	Author string
	AuthorRef string
	Category string
}

type Order struct {
	gorm.Model
	Courses []Course `gorm:"many2many:order_courses;"`
	User string
	TotalPrice float32
}

type CourseWithRating struct {
	Course Course
	Rating float32
}

type CoursesRes struct {
	Courses []Course
	Count int64
}

type Repository interface {
	GetCourses(ctx context.Context, page int) (CoursesRes, error)
	GetCoursesByCategory(ctx context.Context, category string) ([]Course, error)
	GetCourseById(ctx context.Context, id uint) (Course, error)
	CreateCourse(ctx context.Context, course Course) error
	UpdateCourse(ctx context.Context, id uint, title string, description string, price float32, image string, link string, authorRef string) error
	DeleteCourse(ctx context.Context, id uint, authorRef string) error
	CreateOrder(ctx context.Context, order Order) error
	GetUserOrders(ctx context.Context, email string) ([]Order, error)
	GetTeacherCourses(ctx context.Context, email string) ([]Course, error)
}