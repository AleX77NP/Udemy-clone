package user

import(
	"github.com/dgrijalva/jwt-go"
	"time"
	"net/http"
	"strings"
)

var jwtSecret = []byte("myjwtsupersecret")

type Claims struct {
	Email string `json:"email"`
	Role string `json:"role"`
	jwt.StandardClaims
}


func generateToken(email string, role string) (string, error) {
	expTime := time.Now().Add(time.Minute * 1440)
	claims := &Claims{
		Email: email,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func checkToken(w http.ResponseWriter, r *http.Request, n http.Handler ) {
	reqToken := r.Header.Get("Authorization")
	splitJwt := strings.Split(reqToken, "Bearer")

	if len(splitJwt) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqToken = strings.TrimSpace(splitJwt[1])

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(reqToken, claims, func(t *jwt.Token) (interface{},error){
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

func getUser(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	splitJwt := strings.Split(reqToken, "Bearer")

	if len(splitJwt) < 2 {
		return "", InvalidData
	}

	reqToken = strings.TrimSpace(splitJwt[1])

	claims := &Claims{}

	_, err := jwt.ParseWithClaims(reqToken, claims, func(t *jwt.Token) (interface{},error){
		return jwtSecret, nil
	})
	if err != nil {
		return "", InvalidData
	}
	return claims.Email, nil
}