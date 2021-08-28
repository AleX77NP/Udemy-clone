package course

import(
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"strconv"
)

type (
	CreateGetCoursesRequest struct {
		Page int `json:"page"`
	}
	CreateGetCoursesResponse struct {
		Courses []Course `json:"courses"`
		Count int64 `json:"count"`
	}

	CreateGetCoursesByCategoryRequest struct {
		Category string `json:"category"`
	}
	CreateGetCoursesByCategoryResponse struct {
		Courses []Course `json:"courses"`
	}
	CreateGetCourseByIdRequest struct {
		Id uint `json:"id"`
	}
	CreateGetCourseByIdResponse struct {
		CourseR CourseWithRating `json:"course"`
	}
	CreateCourseRequestBody struct {
		Title string `json:"title"`
		Description string `json:"description"`
		Price float32 `json:"price"`
		Image string `json:"image"`
		Link string `json:"link"`
		Author string `json:"author"`
		Category string `json:"category"`
	}
	CreateCourseRequest struct {
		Title string `json:"title"`
		Description string `json:"description"`
		Price float32 `json:"price"`
		Image string `json:"image"`
		Link string `json:"link"`
		Author string `json:"author"`
		AuthorRef string `json:"authorRef"`
		Category string `json:"category"`
	}
	CreateCourseResponse struct {
		Ok string `json:"ok"`
	}
	CreateUpdateCourseRequestBody struct {
		Title string `json:"title"`
		Description string `json:"description"`
		Price float32 `json:"price"`
		Image string `json:"image"`
		Link string `json:"link"`
	}
	CreateUpdateCourseRequest struct {
		Id uint `json:"id"`
		Title string `json:"title"`
		Description string `json:"description"`
		Price float32 `json:"price"`
		Image string `json:"image"`
		Link string `json:"link"`
		AuthorRef string `json:"autorRef"`
	}
	CreateUpdateCourseResponse struct {
		Ok string `json:"ok"`
	}
	CreateDeleteCourseRequest struct {
		Id uint `json:"id"`
		AuthorRef string `json:"autorRef"`
	}
	CreateDeleteCourseResponse struct {
		Ok string `json:"ok"`
	}
	CreateTeacherCoursesRequest struct {
		Email string `json:"email"`
	}
	CreateTeacherCoursesResponse struct {
		Courses []Course `json:"courses"`
	}
	CreateCreateOrderRequestBody struct {
		Courses []Course `json:"courses"`
		TotalPrice float32 `json:"totalPrice"`
	}
	CreateCreateOrderRequest struct {
		Courses []Course `json:"courses"`
		User string `json:"user"`
		TotalPrice float32 `json:"totalPrice"`
	}
	CreateCreateOrderResponse struct {
		Ok string `json:"ok"`
	}
	CreateGetUserOrdersRequest struct {
		Email string `json:"email"`
	}
	CreateGetUserOrdersResponse struct {
		Orders []Order `json:"orders"`
	}
	CreateGetTeacherCoursesRequest struct {
		Email string `json:"email"`
	}
	CreateGetTeacherCoursesResponse struct {
		Courses []Course `json:"courses"`
	}
	CreateRateCourseRequest struct {
		Id uint `json:"id"`
		Rating int32 `json:"rating"`
		User string `json:"user"`
	}
	CreateRateCourseRequestBody struct {
		Rating int32 `json:"rating"`
	}
	CreateRateCourseResponse struct {
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
	case InvalidData:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func stringToUint(id string) (uint, error) {
	u64, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        return 0,err
    }
    wd := uint(u64)
	return wd,nil
}

func decodeGetCoursesReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateGetCoursesRequest
	var page int
	v := r.URL.Query()
	i := v.Get("page")
	p, err := strconv.Atoi(i)
	if err != nil {
		page = 1
	} else {
		page = p
	}
	
	req = CreateGetCoursesRequest{
		Page: page,
	}
	return req, nil
}

func decodeGetCoursesByCategorysReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateGetCoursesByCategoryRequest
	vars := mux.Vars(r)

	req = CreateGetCoursesByCategoryRequest{
		Category: vars["category"],
	}
	return req, nil
}

func decodeGetCourseByIdReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateGetCourseByIdRequest
	vars := mux.Vars(r)

	id := vars["id"]

	uid, err := stringToUint(id)
	if err != nil {
		return nil,err
	}
	req = CreateGetCourseByIdRequest{
		Id: uid,
	}
	return req, nil
}

func decodeCreateCourseReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var reqb CreateCourseRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqb)

	if err != nil {
		return nil, err
	}

	authorRef, err2 := getUser(r)
	if err2 != nil {
		return nil, err2
	}

	req := CreateCourseRequest{
		Title: reqb.Title,
		Description: reqb.Description,
		Price: reqb.Price,
		Image: reqb.Image,
		Link: reqb.Link,
		Author: reqb.Author,
		AuthorRef: authorRef,
		Category: reqb.Category,
	}

	return req,nil
}


func decodeUpdateCourseReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var reqb CreateUpdateCourseRequestBody
	var req CreateUpdateCourseRequest
	err := json.NewDecoder(r.Body).Decode(&reqb)
	
	if err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	id := vars["id"]

	author, err := getUser(r)

	if err != nil {
		return nil, err
	}

	uid, err := stringToUint(id)
	if err != nil {
		return nil,err
	}
	req = CreateUpdateCourseRequest{
		Id: uid,
		Title: reqb.Title,
		Description: reqb.Description,
		Price: reqb.Price,
		Image: reqb.Image,
		Link: reqb.Link,
		AuthorRef: author,
	}
	return req, nil
}

func decodeDeleteCourseReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateDeleteCourseRequest

	author, err := getUser(r)

	if err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	id := vars["id"]

	uid, err := stringToUint(id)
	if err != nil {
		return nil,err
	}
	req = CreateDeleteCourseRequest{
		Id: uid,
		AuthorRef: author,
	}
	return req, nil
}

func decodeCreateOrderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var reqb CreateCreateOrderRequestBody

	user, err := getUser(r)
	if err != nil {
		return nil, err
	}
	err2 := json.NewDecoder(r.Body).Decode(&reqb)
	
	if err2 != nil {
		return nil, err2
	}

	req := CreateCreateOrderRequest{
		Courses: reqb.Courses,
		User: user,
		TotalPrice: reqb.TotalPrice,
	}
	return req, nil
}

func decodeGetUserOrdersRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	user, err := getUser(r)
	if err != nil {
		return nil, err
	}
	req := CreateGetUserOrdersRequest{
		Email: user,
	}
	return req, nil
}

func decodeGetTeacherCoursesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	user, err := getUser(r)
	if err != nil {
		return nil, err
	}
	req := CreateGetTeacherCoursesRequest{
		Email: user,
	}
	return req, nil
}

func decodeRateCourseReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateRateCourseRequest
	var reqb CreateRateCourseRequestBody
	vars := mux.Vars(r)
	id := vars["id"]

	uid, err := stringToUint(id)
	if err != nil {
		return nil,err
	}

	err2 := json.NewDecoder(r.Body).Decode(&reqb)
	
	if err2 != nil {
		return nil, err
	}

	user, err3 := getUser(r)
	if err != nil {
		return nil, err3
	}

	req = CreateRateCourseRequest{
		Id: uid,
		Rating: reqb.Rating,
		User: user,
	}

	return req, nil
}