package course

import (
	"github.com/go-kit/kit/endpoint"

	"context"
)

type Endpoints struct {
	GetCourses endpoint.Endpoint
	GetCoursesByCategory endpoint.Endpoint
	GetCourseById endpoint.Endpoint
	CreateCourse endpoint.Endpoint
	UpdateCourse endpoint.Endpoint
	DeleteCourse endpoint.Endpoint
	CreateOrder endpoint.Endpoint
	GetUserOrders endpoint.Endpoint
	GetTeacherCourses endpoint.Endpoint
	RateCourse endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetCourses: makeGetCoursesEndpoint(s),
		GetCoursesByCategory: makeGetCoursesByCategoryEndpoint(s),
		GetCourseById: makeGetCourseByIdEndpoint(s),
		CreateCourse: makeCreateCourseEndpoint(s),
		UpdateCourse: makeUpdateCourseEndpoint(s),
		DeleteCourse: makeDeleteCourseEndpoint(s),
		CreateOrder: makeCreateOrderEndpoints(s),
		GetUserOrders: makeGetUserOrdersEndpoints(s),
		GetTeacherCourses: makeGetTeacherCoursesEndpoints(s),
		RateCourse: makeRateCourseEndpoint(s),
	}
}

func makeGetCoursesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateGetCoursesRequest)
		rc, err := s.GetCourses(ctx, req.Page)
		return CreateGetCoursesResponse{Courses: rc.Courses, Count: rc.Count}, err
	}
}

func makeGetCoursesByCategoryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateGetCoursesByCategoryRequest)
		rc, err := s.GetCoursesByCategory(ctx, req.Category)
		return CreateGetCoursesByCategoryResponse{Courses: rc}, err
	}
}

func makeGetCourseByIdEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateGetCourseByIdRequest)
		rc, err := s.GetCourseById(ctx, req.Id)
		return CreateGetCourseByIdResponse{CourseR: rc}, err
	}
}

func makeCreateCourseEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateCourseRequest)
		rc, err := s.CreateCourse(ctx, req.Title, req.Description, req.Price, req.Image, req.Link, req.Author, req.AuthorRef, req.Category)
		return CreateCourseResponse{Ok: rc}, err
	}
}

func makeRateCourseEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRateCourseRequest)
		rc, err := s.RateCourse(ctx, req.Id, req.Rating, req.User)
		return CreateRateCourseResponse{Ok: rc}, err
	}
}

func makeUpdateCourseEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUpdateCourseRequest)
		rc, err := s.UpdateCourse(ctx, req.Id, req.Title, req.Description, req.Price, req.Image, req.Link, req.AuthorRef)
		return CreateCourseResponse{Ok: rc}, err
	}
}

func makeDeleteCourseEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateDeleteCourseRequest)
		rc, err := s.DeleteCourse(ctx, req.Id, req.AuthorRef)
		return CreateCourseResponse{Ok: rc}, err
	}
}

func makeCreateOrderEndpoints(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateCreateOrderRequest)
		rc, err := s.CreateOrder(ctx, req.Courses, req.User, req.TotalPrice)
		return CreateCreateOrderResponse{Ok: rc}, err
	}
}

func makeGetUserOrdersEndpoints(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateGetUserOrdersRequest)
		rc, err := s.GetUserOrders(ctx, req.Email)
		return CreateGetUserOrdersResponse{Orders: rc}, err
	}
}

func makeGetTeacherCoursesEndpoints(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateGetTeacherCoursesRequest)
		rc, err := s.GetTeacherCourses(ctx, req.Email)
		return CreateGetTeacherCoursesResponse{Courses: rc}, err
	}
}