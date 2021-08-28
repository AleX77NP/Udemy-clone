package server

import (
	"github.com/hashicorp/go-hclog"
	"context"
	protos "lurushop.com/rating-service/protos/ratings"
	//"github.com/go-kit/log"
)


var (
	RRepo Repository
	Srv Service
)


type RatingService struct {
	log hclog.Logger
}

func NewRatingService(l hclog.Logger) *RatingService {
	return &RatingService{
		log: l,
	}
}


func(c *RatingService) GetRating(ctx context.Context, rr *protos.GetRatingRequest) (*protos.GetRatingResponse, error) {
	c.log.Info("Handle Get Rating", "id", rr.GetId())

	rt, err := Srv.GetCourseRating(ctx, rr.GetId())
	if err != nil {
		return nil, err
	}

	return &protos.GetRatingResponse{
		Rating: rt,
	}, nil
}

func(c *RatingService) AddRating(ctx context.Context, rr *protos.AddRatingRequest) (*protos.AddRatingResponse, error) {
	c.log.Info("Handle Add Rating", "id", rr.GetId(), "rating", rr.GetRating(), "user", rr.GetUser())

	if _, err := Srv.AddCourseRating(ctx, rr.GetId(), rr.GetRating(), rr.GetUser()); err != nil {
		return nil, err
	}

	return &protos.AddRatingResponse{
		Ok: "ok",
	}, nil
}