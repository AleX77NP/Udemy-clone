package course

import (
	"context"

	"github.com/go-kit/log"

	"errors"

	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)

type repo struct {
	db     *gorm.DB
	logger log.Logger
}

var RepoErr = errors.New("Could not handle your request.")

func NewRepo(db *gorm.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (repo *repo) GetCourses(ctx context.Context, page int) (CoursesRes, error) {
	db := repo.db

	// offset := (page-1) * 6

	var courses []Course
	var cnt int64

	if result := db.Model(&Course{}).Count(&cnt); result.Error != nil {
		return CoursesRes{}, result.Error
	}

	if result2 := db.Find(&courses); result2.Error != nil {
		return CoursesRes{}, result2.Error
	}

	res := CoursesRes{
		Courses: courses,
		Count:   cnt,
	}

	return res, nil
}

func (repo *repo) GetCoursesByCategory(ctx context.Context, category string) ([]Course, error) {
	db := repo.db

	var courses []Course
	if result := db.Where("category = ?", category).Find(&courses); result.Error != nil {
		return courses, result.Error
	}

	return courses, nil
}

func (repo *repo) GetCourseById(ctx context.Context, id uint) (Course, error) {
	db := repo.db

	var course Course
	if result := db.First(&course, id); result.Error != nil {
		return course, NotFound
	}

	return course, nil
}

func (repo *repo) CreateCourse(ctx context.Context, course Course) error {
	db := repo.db

	if course.Title == "" || course.Description == "" || course.Author == "" || course.AuthorRef == "" || course.Image == "" || course.Link == "" || course.Price < 0 {
		return InvalidData
	}

	if err := db.Create(&course).Error; err != nil {
		return err
	}

	return nil
}

func (repo *repo) UpdateCourse(ctx context.Context, id uint, title string, description string, price float32, image string, link string, authorRef string) error {
	db := repo.db

	var course Course
	db.First(&course, id)

	if course.AuthorRef != authorRef {
		return Unauthorized
	}

	if err := db.Model(&course).Updates(Course{Title: title, Description: description, Price: price, Image: image, Link: link}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *repo) DeleteCourse(ctx context.Context, id uint, authorRef string) error {
	db := repo.db

	var course Course
	db.First(&course, id)

	if course.AuthorRef != authorRef {
		return Unauthorized
	}

	if err := db.Delete(&Course{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *repo) GetTeacherCourses(ctx context.Context, email string) ([]Course, error) {
	db := repo.db

	var courses []Course
	if err := db.Where("author_ref = ?", email).Find(&courses).Error; err != nil {
		return []Course{}, RepoErr
	}
	return courses, nil
}

func (repo *repo) CreateOrder(ctx context.Context, order Order) error {
	db := repo.db

	if order.TotalPrice <= 0 || order.User == "" {
		return InvalidData
	}

	if err := db.Create(&order).Error; err != nil {
		return err
	}
	return nil
}

func (repo *repo) GetUserOrders(ctx context.Context, email string) ([]Order, error) {
	db := repo.db

	var orders []Order
	if email == "" {
		return []Order{}, InvalidData
	}
	if err := db.Where("user = ?", email).Preload("Courses").Find(&orders).Error; err != nil {
		return []Order{}, RepoErr
	}
	return orders, nil
}
