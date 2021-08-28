package user

import (
	"context"
	"net/http"

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

	ar := r.PathPrefix("/api/users/private").Subrouter()
	ar.Use(jwtMiddleware)

	ar.Methods("GET", "OPTIONS").Path("/me").Handler(httptransport.NewServer(
		endpoints.GetUser,
		decodeGetUserReq,
		encodeResponse,
		options...,
	))

	ar.Methods("PUT", "OPTIONS").Path("/change").Handler(httptransport.NewServer(
		endpoints.ChangeRole,
		decodeChangeRoleReq,
		encodeResponse,
		options...,
	))

	ar.Methods("PUT", "OPTIONS").Path("/password/reset").Handler(httptransport.NewServer(
		endpoints.ChangePassword,
		decodeChangePasswordReq,
		encodeResponse,
		options...,
	))

	ar.Methods("OPTIONS", "PUT").Path("/profile/edit").Handler(httptransport.NewServer(
		endpoints.EditProfile,
		decodeEditProfileReq,
		encodeResponse,
		options...,
	))

	ar.Use(CORS)

	r.Methods("POST").Path("/api/users/signup").Handler(httptransport.NewServer(
		endpoints.CreateUser,
		decodeUserReq,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/api/users/login").Handler(httptransport.NewServer(
		endpoints.LoginUser,
		decodeLoginReq,
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

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		checkToken(w, r, next)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
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
