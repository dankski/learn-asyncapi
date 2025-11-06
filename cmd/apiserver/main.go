package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/dankski/learn-asyncapi/apiserver"
	"github.com/dankski/learn-asyncapi/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conf, err := config.New()
	if err != nil {
		return err
	}

	server := apiserver.New(conf)
	if err := server.Start(ctx); err != nil {
		return err
	}

	return nil
}
