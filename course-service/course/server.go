package course

import (
	"context"
	"net/http"

	//"fmt"

	httptransport "github.com/go-kit/kit/transport/http"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints, options []kithttp.ServerOption) http.Handler {
	var (
		r            = mux.NewRouter()
		errorEncoder = kithttp.ServerErrorEncoder(encodeErrorResponse)
	)
	options = append(options, errorEncoder)

	r.Use(commonMiddleware)
	r.Use(CORS)

	ar := r.PathPrefix("/api/courses/student").Subrouter()
	ar.Use(authMiddleware)

	ar.Methods("POST", "OPTIONS").Path("/order").Handler(httptransport.NewServer(
		endpoints.CreateOrder,
		decodeCreateOrderRequest,
		encodeResponse,
		options...,
	))

	ar.Methods("POST", "OPTIONS").Path("/rate/{id}").Handler(httptransport.NewServer(
		endpoints.RateCourse,
		decodeRateCourseReq,
		encodeResponse,
		options...,
	))

	ar.Methods("GET", "OPTIONS").Path("/orders/me").Handler(httptransport.NewServer(
		endpoints.GetUserOrders,
		decodeGetUserOrdersRequest,
		encodeResponse,
		options...,
	))

	ar.Use(CORS)

	tr := r.PathPrefix("/api/courses/teacher").Subrouter()
	tr.Use(authTeacherMiddleware)

	tr.Methods("POST", "OPTIONS").Path("/create").Handler(httptransport.NewServer(
		endpoints.CreateCourse,
		decodeCreateCourseReq,
		encodeResponse,
		options...,
	))

	tr.Methods("PUT", "OPTIONS").Path("/update/{id}").Handler(httptransport.NewServer(
		endpoints.UpdateCourse,
		decodeUpdateCourseReq,
		encodeResponse,
		options...,
	))

	tr.Methods("DELETE", "OPTIONS").Path("/delete/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteCourse,
		decodeDeleteCourseReq,
		encodeResponse,
		options...,
	))

	tr.Methods("GET", "OPTIONS").Path("/my").Handler(httptransport.NewServer(
		endpoints.GetTeacherCourses,
		decodeGetTeacherCoursesRequest,
		encodeResponse,
		options...,
	))

	tr.Use(CORS)

	r.Methods("GET").Path("/api/courses").Handler(httptransport.NewServer(
		endpoints.GetCourses,
		decodeGetCoursesReq,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/api/courses/category/{category}").Handler(httptransport.NewServer(
		endpoints.GetCoursesByCategory,
		decodeGetCoursesByCategorysReq,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/api/courses/{id}").Handler(httptransport.NewServer(
		endpoints.GetCourseById,
		decodeGetCourseByIdReq,
		encodeResponse,
		options...,
	))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkToken(w, r, next)
	})
}

func authTeacherMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkTokenTeacher(w, r, next)
	})
}

func authAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkTokenAdmin(w, r, next)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Next
		next.ServeHTTP(w, r)
		return
	})
}
