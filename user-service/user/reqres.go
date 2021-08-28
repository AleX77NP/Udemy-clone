package user

import (
	"context"
	"encoding/json"
	"net/http"
	//"github.com/gorilla/mux"
)

type (
	CreateUserRequest struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Name      string `json:"name"`
		Surname   string `json:"surname"`
		Role      string `json:"role"`
		Image     string `json:"image"`
		ResetCode string `json:"resetCode"`
	}
	CreateUserResponse struct {
		ResetCode string `json:"resetCode"`
	}
	CreateLoginUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	CreateLoginUserResponse struct {
		Token string `json:"token"`
	}
	CreateGetUserRequest struct {
		Email string `json:"email"`
	}
	CreateGetUserResponse struct {
		UserProfile UserProfile `json:"userProfile"`
	}
	CreateChangePasswordRequestBody struct {
		ResetCode   string `json:"resetCode"`
		NewPassword string `json:"newPassword"`
	}
	CreateChangePasswordRequest struct {
		Email       string `json:"email"`
		ResetCode   string `json:"resetCode"`
		NewPassword string `json:"newPassword"`
	}
	CreateChangePasswordResponse struct {
		Ok string `json:"ok"`
	}
	CreateEditProfileRequestBody struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Image   string `json:"image"`
	}
	CreateEditProfileRequest struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Image   string `json:"image"`
		Role    string `json:"role"`
	}
	CreateEditProfileResponse struct {
		Ok string `json:"ok"`
	}
	CreateChangeRoleRequestBody struct {
		Email string `json:"email"`
	}
	CreateChangeRoleRequest struct {
		Email string `json:"email`
	}
	CreateChangeRoleResponse struct {
		Ok string `json:"ok"`
	}
)

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeErrorResponse(ctx context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encoded error with nil error.")
	}
	w.WriteHeader(codeFromErr(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFromErr(err error) int {
	switch err {
	case NotFound:
		return http.StatusNotFound
	case WrongPassword, UniqueEmail, InvalidData, WrongCode:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func decodeUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeLoginReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateLoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	user, err := getUser(r)
	if err != nil {
		return nil, InvalidData
	}
	req := CreateGetUserRequest{
		Email: user,
	}
	return req, nil
}

func decodeChangePasswordReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var reqb CreateChangePasswordRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqb)

	if err != nil {
		return nil, InvalidData
	}
	user, err := getUser(r)
	if err != nil {
		return nil, InvalidData
	}

	req := CreateChangePasswordRequest{
		Email:       user,
		ResetCode:   reqb.ResetCode,
		NewPassword: reqb.NewPassword,
	}
	return req, nil
}

func decodeEditProfileReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var reqb CreateEditProfileRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqb)

	if err != nil {
		return nil, InvalidData
	}
	user, err := getUser(r)
	if err != nil {
		return nil, InvalidData
	}

	req := CreateEditProfileRequest{
		Email:   user,
		Name:    reqb.Name,
		Surname: reqb.Surname,
		Image:   reqb.Image,
	}
	return req, nil
}

func decodeChangeRoleReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var reqb CreateChangeRoleRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqb)

	if err != nil {
		return nil, InvalidData
	}
	req := CreateChangeRoleRequest{
		Email: reqb.Email,
	}
	return req, nil
}
