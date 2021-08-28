package server

import(
	"github.com/go-kit/log"
	"context"

	"gorm.io/gorm"
	"errors"
)

var (
	DBConn *gorm.DB
)

type repo struct {
	db *gorm.DB
	logger log.Logger
}

var RepoErr = errors.New("Could not handle your request.")

func NewRepo(db *gorm.DB, logger log.Logger) Repository {
	return &repo{
		db: db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (repo *repo) GetCourseRating(ctx context.Context, courseId uint32) (float32, error) {
	db := repo.db
	var r float32

	if err := db.Model(&CourseRating{}).Select("avg(rating) as r").Where("course_id = ?", courseId).Group("course_id").First(&r).Error; err != nil {
		return 0,NotFound
	}

	return r,nil
}

func (repo *repo) AddCourseRating(ctx context.Context, courseRating CourseRating) error {
	db := repo.db

	var alreadyRated CourseRating
	if err := db.Where("course_id = ? AND user = ?", courseRating.CourseId, courseRating.User).First(&alreadyRated).Error; err != nil {
		if courseRating.Rating < 1 || courseRating.User == "" {
			return InvalidData
		}
		if err := db.Create(&courseRating).Error; err != nil {
			return RepoErr
		}
		return nil
	} else {
		if err := db.Model(CourseRating{}).Where("course_id = ? AND user = ?", courseRating.CourseId, courseRating.User).Update("rating", courseRating.Rating).Error; err != nil {
			return RepoErr
		}
		return nil
	}
}