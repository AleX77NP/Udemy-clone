module lurushop.com/course-service

go 1.16

replace lurushop.com/rating-service => ../rating-service

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-kit/kit v0.11.0
	github.com/go-kit/log v0.1.0
	github.com/gorilla/mux v1.8.0
	google.golang.org/grpc v1.39.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
	lurushop.com/rating-service v0.0.0-00010101000000-000000000000
)
