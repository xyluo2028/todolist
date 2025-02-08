package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todolist/internal/handlers"
	"todolist/internal/middleware"
	"todolist/internal/store"
)

func main() {

	fmt.Println("Welcome to the Todo List!")

	taskStore := store.NewTaskStore()
	userStore := store.NewUserStore()

	taskHandler := handlers.NewTaskHandler(taskStore)
	welcomeHandler := handlers.NewWelcomeHandler()
	userHandler := handlers.NewUserHandler(userStore, taskStore)
	auth := middleware.NewAuthMiddleware(userStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/welcome", auth.Authenticate(welcomeHandler.Welcome))
	mux.HandleFunc("/printTasks", auth.Authenticate(taskHandler.PrintTasksHttp))
	mux.HandleFunc("/addTask", auth.Authenticate(taskHandler.AddTaskHttp))
	mux.HandleFunc("/removeTask", auth.Authenticate(taskHandler.RemoveTaskHttp))
	mux.HandleFunc("/updateTask", auth.Authenticate(taskHandler.UpdateTaskHttp))
	mux.HandleFunc("/deactivate", auth.Authenticate(userHandler.DeleteUser))
	mux.HandleFunc("/register", userHandler.Register)

	server := &http.Server{
		Addr:    ":7071",
		Handler: mux,
	}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		log.Printf("shutting down server..")
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
