package course

import (
	protos "lurushop.com/rating-service/protos/ratings"
	"context"
	"fmt"
)

var (
	RatingClient protos.RatingClient
)

func GetCourseRating(ctx context.Context, id uint32) (float32, error) {
	rr := &protos.GetRatingRequest{
		Id: id,
	}
	res, err := RatingClient.GetRating(ctx, rr)
	if err != nil {
		fmt.Println(err)
		return 0,GetRatingError
	}
	return res.Rating, nil
}

func AddCourseRating(ctx context.Context, id uint32, rating int32, user string) (string, error) {
	ar := &protos.AddRatingRequest{
		Id: id,
		Rating: rating,
		User: user,
	}
	res, err := RatingClient.AddRating(ctx, ar)
	if err != nil {
		fmt.Println(err)
		return "", AddRatingError
	}
	return res.Ok, nil
}