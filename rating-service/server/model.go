package server

import (
	"gorm.io/gorm"
	"context"
)

type CourseRating struct {
	gorm.Model
	CourseId uint32
	Rating int32
	User string
}

type Repository interface {
	GetCourseRating(ctx context.Context, productId uint32) (float32, error)
	AddCourseRating(ctx context.Context, courseRating CourseRating) error
}