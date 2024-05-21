package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/mikitabablo/webpages/config"
	"github.com/mikitabablo/webpages/internal/app"
	"github.com/rs/zerolog"
)

func main() {
	cfg, err := config.Load("config.toml")
	if err != nil {
		log.Panicf("cannot load configs: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// setting up log
	l, err := loadLogger(cfg.Application.LogLevel)
	if err != nil {
		log.Panicf("cannot load logger: %s", err)
	}

	a := app.NewApp(cfg, l)
	a.Run(ctx, wg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	cancel()
	wg.Wait()
}

func loadLogger(logLevel string) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		return zerolog.Logger{}, err
	}

	zerolog.SetGlobalLevel(level)
	return zerolog.New(zerolog.MultiLevelWriter(os.Stdout)).
		With().
		Caller().
		Timestamp().
		Logger(), nil
}
