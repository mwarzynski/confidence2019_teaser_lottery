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
	DefaultFlag   = "confidence{flag}"
	DefaultListen = ":8080"

	DefaultLotteryPeriod      = time.Minute * 12
	DefaultAccountRemoveAfter = time.Minute * 15
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse environment variables.
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = DefaultListen
	}

	flag := os.Getenv("FLAG")
	if flag == "" {
		flag = DefaultFlag
	}

	lotteryPeriod := DefaultLotteryPeriod
	if d, err := time.ParseDuration(os.Getenv("LOTTERY_PERIOD")); err == nil {
		lotteryPeriod = d
	}
	deleteAccountAfter := DefaultAccountRemoveAfter
	if d, err := time.ParseDuration(os.Getenv("DELETE_ACCOUNT_AFTER")); err == nil {
		deleteAccountAfter = d
	}

	service := app.NewService(ctx, lotteryPeriod, deleteAccountAfter)
	router := transport.InitRouter(service, flag)

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
