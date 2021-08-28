package server

import (

)

import (
	"context"
	"errors"
)

type Service interface {
	GetCourseRating(ctx context.Context, courseId uint32) (float32, error)
	AddCourseRating(ctx context.Context, courseId uint32, rating int32, user string) (string, error)
}

var (
	NotFound = errors.New("Course not found.")
	InvalidData = errors.New("Invalid data")
)