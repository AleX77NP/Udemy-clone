package server

import (
	"github.com/go-kit/log/level"
	"github.com/go-kit/log"
	"context"
)

type service struct {
	repository Repository
	logger log.Logger
}

func NewService(rep Repository, logger log.Logger) Service {
	return &service{
		repository: rep,
		logger: logger,
	}
}

func (s service) GetCourseRating(ctx context.Context, courseId uint32) (float32, error) {
	logger := log.With(s.logger, "method", "GetCourseRating")

	rating, err := s.repository.GetCourseRating(ctx, courseId)
	if err != nil {
		level.Error(logger).Log("err",err)
		return 0, err
	}
	logger.Log("success getting rating")
	return rating, nil
}

func (s service) AddCourseRating(ctx context.Context, courseId uint32, rating int32, user string) (string, error) {
	logger := log.With(s.logger, "method", "GetCourseRating")
	cr := CourseRating{
		CourseId: courseId,
		Rating: rating,
		User: user,
	}

	if err := s.repository.AddCourseRating(ctx, cr); err != nil {
		level.Error(logger).Log("err",err)
		return "", err
	}

	logger.Log("success adding rating")
	return "Ok", nil
}