package main

import (
	mainWire "flowforge-api/cmd/wire"
	"flowforge-api/infrastructure/config"
	"flowforge-api/infrastructure/messaging/rabbitmq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	env := config.Load()
	db := config.ConnectDatabase(env)

	container := mainWire.NewContainer(db, env)

	consumer := rabbitmq.NewWorkerConsumer(
		env,
		container.ConsumerHandler,
	)

	go func() {
		if err := consumer.Start(); err != nil {
			log.Fatalf("failed to start consumer: %v", err)
		}
	}()

	log.Println("worker started")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh

	log.Printf("received signal %s, shutting down", sig)

	if err := consumer.Stop(); err != nil {
		log.Printf("failed to stop consumer: %v", err)
	}
}
