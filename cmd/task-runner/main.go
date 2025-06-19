package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kylerqws/task-runner/internal/bootstrap"
	"github.com/kylerqws/task-runner/internal/domain/service"
	"github.com/kylerqws/task-runner/internal/transport/http/handler"
	"github.com/kylerqws/task-runner/internal/transport/http/router"
)

const serverAddr = ":8080"

func main() {
	manager := initManager()
	server := initServer(manager)

	waitForShutdown(server)
}

func initManager() *service.TaskManager {
	manager := service.NewTaskManager()
	bootstrap.RegisterTaskFactories(manager)

	return manager
}

func initServer(manager *service.TaskManager) *http.Server {
	taskHandler := handler.NewTaskHandler(manager)
	httpHandler := router.InitTaskRouter(taskHandler)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: httpHandler,
	}

	go func() {
		log.Println("Server listening on", serverAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	return server
}

func waitForShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
