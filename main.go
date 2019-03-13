package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mwarzynski/confidence_web/app"
	"github.com/mwarzynski/confidence_web/transport"
)

var (
	DefaultListen        = ":8080"
	DefaultLotteryPeriod = time.Minute * 10
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse environment variables.
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = DefaultListen
	}

	lotteryPeriod := DefaultLotteryPeriod
	if d, err := time.ParseDuration(os.Getenv("LOTTERY_PERIOD")); err == nil {
		lotteryPeriod = d
	}

	service := app.NewService(ctx, lotteryPeriod)
	router := transport.InitRouter(service)

	timeout := time.Duration(time.Minute)
	server := &http.Server{
		Addr:         listen,
		Handler:      router,
		ReadTimeout:  time.Second * time.Duration(timeout),
		WriteTimeout: time.Second * time.Duration(timeout),
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("server shutdown: %s", err)
		}
	}()
	server.ListenAndServe()
}
