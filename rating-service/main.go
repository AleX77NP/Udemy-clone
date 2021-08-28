package main

import (
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	protos "lurushop.com/rating-service/protos/ratings"
	"lurushop.com/rating-service/server"
	"net"
	"os"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"github.com/go-kit/log"
)

var err error

func main() {

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}


	server.RRepo = server.NewRepo(server.DBConn, logger)
	server.Srv = server.NewService(server.RRepo, logger)

	log := hclog.Default()
	gs := grpc.NewServer()
	rs := server.NewRatingService(log)

	protos.RegisterRatingServer(gs, rs)

	reflection.Register(gs)

	server.DBConn, err = gorm.Open(sqlite.Open("ratings.db"), &gorm.Config{})
	if err != nil {
		panic("DB connection error...")
	}

	server.RRepo = server.NewRepo(server.DBConn, logger)
	server.Srv = server.NewService(server.RRepo, logger)

	server.DBConn.AutoMigrate(&server.CourseRating{})

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	gs.Serve(l)
}