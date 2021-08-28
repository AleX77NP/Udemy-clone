package course

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type service struct {
	repository Repository
	logger     log.Logger
}

func NewService(rep Repository, logger log.Logger) Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

func (s service) GetCourses(ctx context.Context, page int) (CoursesRes, error) {
	logger := log.With(s.logger, "method", "GetCourses")

	if courses, err := s.repository.GetCourses(ctx, page); err != nil {
		level.Error(logger).Log("err", err)
		return CoursesRes{}, err
	} else {
		logger.Log("Get courses")
		return courses, nil
	}
}

func (s service) GetCoursesByCategory(ctx context.Context, category string) ([]Course, error) {
	logger := log.With(s.logger, "method", "GetCoursesByCategory")

	if courses, err := s.repository.GetCoursesByCategory(ctx, category); err != nil {
		level.Error(logger).Log("err", err)
		return []Course{}, err
	} else {
		logger.Log("Get courses by category")
		return courses, nil
	}
}

func (s service) GetCourseById(ctx context.Context, id uint) (CourseWithRating, error) {
	logger := log.With(s.logger, "method", "GetCourseById")
	var courseR CourseWithRating
	var r float32

	u := uint32(id)

	if course, err := s.repository.GetCourseById(ctx, id); err != nil {
		return CourseWithRating{}, err
	} else {
		res, err := GetCourseRating(ctx, u)
		if err != nil {
			r = 0
		} else {
			r = res
		}

		courseR = CourseWithRating{
			Course: course,
			Rating: r,
		}
		logger.Log("Get course by id")
		return courseR, nil
	}
}

func (s service) CreateCourse(ctx context.Context, title string, description string, price float32, image string, link string, author string, authorRef string, category string) (string, error) {
	logger := log.With(s.logger, "method", "CreateCourse")

	course := Course{
		Title:       title,
		Description: description,
		Price:       price,
		Image:       image,
		Link:        link,
		Author:      author,
		AuthorRef:   authorRef,
		Category:    category,
	}

	if err := s.repository.CreateCourse(ctx, course); err != nil {
		return "", err
	}

	logger.Log("Create course")
	return "Course created.", nil
}

func (s service) RateCourse(ctx context.Context, id uint, rating int32, user string) (string, error) {
	logger := log.With(s.logger, "method", "RateCourse")

	u := uint32(id)
	res, err := AddCourseRating(ctx, u, rating, user)
	if err != nil {
		return "", err
	}

	logger.Log("Rate course")
	return res, nil
}

func (s service) UpdateCourse(ctx context.Context, id uint, title string, description string, price float32, image string, link string, authorRef string) (string, error) {
	logger := log.With(s.logger, "method", "UpdateCourse")

	if err := s.repository.UpdateCourse(ctx, id, title, description, price, image, link, authorRef); err != nil {
		return "", err
	}
	logger.Log("Update course")

	return "Course updated.", nil
}

func (s service) DeleteCourse(ctx context.Context, id uint, authorRef string) (string, error) {
	logger := log.With(s.logger, "method", "DeleteCourse")

	if err := s.repository.DeleteCourse(ctx, id, authorRef); err != nil {
		return "", err
	}
	logger.Log("Delete course")

	return "Course deleted", nil
}

func (s service) CreateOrder(ctx context.Context, courses []Course, user string, totalPrice float32) (string, error) {
	logger := log.With(s.logger, "method", "CreateOrder")

	order := Order{
		Courses:    courses,
		User:       user,
		TotalPrice: totalPrice,
	}
	if err := s.repository.CreateOrder(ctx, order); err != nil {
		return "", err
	}
	logger.Log("Order created")
	return "Order created.", nil
}

func (s service) GetUserOrders(ctx context.Context, email string) ([]Order, error) {
	logger := log.With(s.logger, "method", "GetUserOrders")

	orders, err := s.repository.GetUserOrders(ctx, email)
	if err != nil {
		return orders, err
	}
	logger.Log("User's orders")
	return orders, nil
}

func (s service) GetTeacherCourses(ctx context.Context, email string) ([]Course, error) {
	logger := log.With(s.logger, "method", "GetTeacherCourses")

	courses, err := s.repository.GetTeacherCourses(ctx, email)
	if err != nil {
		return courses, err
	}

	logger.Log("Teacher's courses")
	return courses, nil
}
