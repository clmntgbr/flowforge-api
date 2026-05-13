package main

import (
	mainWire "flowforge-api/cmd/wire"
	"flowforge-api/infrastructure/config"
	"flowforge-api/infrastructure/rabbitmq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	env := config.Load()
	db := config.ConnectDatabase(env)

	container := mainWire.NewContainer(db, env)

	runner := rabbitmq.NewWorkerRunner(
		env,
		container.RunnerHandler,
	)

	go func() {
		if err := runner.Start(); err != nil {
			log.Fatalf("failed to start runner: %v", err)
		}
	}()

	log.Println("runner started")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh

	log.Printf("received signal %s, shutting down", sig)

	if err := runner.Stop(); err != nil {
		log.Printf("failed to stop runner: %v", err)
	}
}
