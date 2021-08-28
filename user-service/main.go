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
	"lurushop.com/user-service/user"
)

var err error


func main() {
	var httpAddr = flag.String("http", ":8081", "http listen address")
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


	flag.Parse()
	ctx := context.Background()

	user.DBConn, err = gorm.Open(sqlite.Open("shop.db"), &gorm.Config{})
	if err != nil {
		panic("DB connection error...")
	}
	user.DBConn.AutoMigrate(&user.User{})


	var srv user.Service
	{
		repository := user.NewRepo(user.DBConn, logger)

		srv = user.NewService(repository, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := user.MakeEndpoints(srv)

	go func() {
		serverOption := []kithttp.ServerOption{}
		fmt.Println("listening on port", *httpAddr)
		handler := user.NewHTTPServer(ctx, endpoints, serverOption)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}