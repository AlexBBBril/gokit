package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/AlexBBBril/gokit/internal/app/gokit/service"
	"github.com/AlexBBBril/gokit/internal/app/gokit/service/cockroachdb"
	httptransport "github.com/AlexBBBril/gokit/internal/app/gokit/transport/http"
	"github.com/AlexBBBril/gokit/internal/app/gokit/transport/order"
	"github.com/AlexBBBril/toolkit/env"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "order",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	checkErr(level.Info(logger).Log("msg", "service started"))
	defer func() {
		checkErr(level.Info(logger).Log("msg", "service ended"))
	}()

	var db *sql.DB
	{
		var err error
		db, err = sql.Open(env.Required("DB_NAME"), env.Required("POSTGRES_URL"))
		if nil != err {
			exitLog(logger, err)
		}
	}

	var svc service.OrderService
	{
		repository, err := cockroachdb.New(db, logger)
		if nil != err {
			exitLog(logger, err)
		}
		svc = service.NewService(repository, logger)
	}

	handler := httptransport.NewService(order.MakeEndpoints(svc), logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		_ = level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: handler,
		}

		errs <- server.ListenAndServe()
	}()

	_ = level.Error(logger).Log("exit", <-errs)
}

func checkErr(err error) {
	if nil != err {
		panic(fmt.Errorf("%w", err))
	}
}

func exitLog(logger log.Logger, err error) {
	_ = level.Error(logger).Log("exit", err)
	os.Exit(-1)
}
