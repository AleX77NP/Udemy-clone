package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/log/level"

	"net/http"
	"os"
	"os/signal"
	"syscall"

    "gorm.io/gorm"
	"gorm.io/driver/sqlite"
	kithttp "github.com/go-kit/kit/transport/http"
	"lurushop.com/course-service/course"
	protos "lurushop.com/rating-service/protos/ratings"

	"google.golang.org/grpc"
)

var err error

func main() {
	var httpAddr = flag.String("http", ":8082", "http listen address")
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

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	gConn, gErr := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if gErr != nil {
		panic(gErr)
	}
	defer gConn.Close()

	//grpc rating client
	course.RatingClient = protos.NewRatingClient(gConn)

	flag.Parse()
	ctx := context.Background()

	course.DBConn, err = gorm.Open(sqlite.Open("shop-courses.db"), &gorm.Config{})
	if err != nil {
		panic("DB connection error...")
	}
	course.DBConn.AutoMigrate(&course.Course{})
	course.DBConn.AutoMigrate(&course.Order{})


	var srv course.Service
	{
		repository := course.NewRepo(course.DBConn, logger)

		srv = course.NewService(repository, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := course.MakeEndpoints(srv)

	go func() {
		serverOption := []kithttp.ServerOption{}
		fmt.Println("listening on port", *httpAddr)
		handler := course.NewHTTPServer(ctx, endpoints, serverOption)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}