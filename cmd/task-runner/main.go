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

// main is the application entry point.
// It initializes the task manager, HTTP server, and handles graceful shutdown.
func main() {
	manager := initManager()
	server := initServer(manager)

	waitForShutdown(server)
}

// initManager creates a new TaskManager and registers all available task factories.
func initManager() *service.TaskManager {
	manager := service.NewTaskManager()
	bootstrap.RegisterTaskFactories(manager)

	return manager
}

// initServer configures and starts the HTTP server with the task routes.
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

// waitForShutdown blocks until a termination signal is received
// and then shuts down the HTTP server gracefully.
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
