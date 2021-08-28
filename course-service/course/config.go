package course

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("myjwtsupersecret")

var (
	STUDENT = "student"
	TEACHER = "teacher"
	ADMIN   = "admin"
)

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func checkToken(w http.ResponseWriter, r *http.Request, n http.Handler) {
	reqToken := r.Header.Get("Authorization")
	splitJwt := strings.Split(reqToken, "Bearer")

	if len(splitJwt) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqToken = strings.TrimSpace(splitJwt[1])

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(reqToken, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	n.ServeHTTP(w, r)
}

func checkTokenTeacher(w http.ResponseWriter, r *http.Request, n http.Handler) {
	reqToken := r.Header.Get("Authorization")
	splitJwt := strings.Split(reqToken, "Bearer")

	if len(splitJwt) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqToken = strings.TrimSpace(splitJwt[1])

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(reqToken, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if claims.Role != TEACHER {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	n.ServeHTTP(w, r)
}

func checkTokenAdmin(w http.ResponseWriter, r *http.Request, n http.Handler) {
	reqToken := r.Header.Get("Authorization")
	splitJwt := strings.Split(reqToken, "Bearer")

	if len(splitJwt) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqToken = strings.TrimSpace(splitJwt[1])

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(reqToken, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if claims.Role != ADMIN {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	n.ServeHTTP(w, r)
}

func getUser(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	splitJwt := strings.Split(reqToken, "Bearer")

	if len(splitJwt) < 2 {
		return "", InvalidData
	}

	reqToken = strings.TrimSpace(splitJwt[1])

	claims := &Claims{}

	_, err := jwt.ParseWithClaims(reqToken, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return "", InvalidData
	}
	return claims.Email, nil
}
