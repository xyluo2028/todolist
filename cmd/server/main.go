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
	"todolist/internal/repository"
	"todolist/internal/services"

	"github.com/gocql/gocql"
)

func main() {

	fmt.Println("Welcome to the Todo List!")

	cluster := gocql.NewCluster("127.0.0.1:9042") // Replace with your Cassandra node IPs
	cluster.Keyspace = "todolist"                 // The keyspace you created
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second // Example timeout

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Could not connect to Cassandra: %v", err)
	}
	defer session.Close()

	//taskRepo := repository.NewInMemTaskRepository()
	//userRepo := repository.NewInMemUserRepository()
	taskRepo := repository.NewCassandraTaskRepository(session)
	userRepo := repository.NewCassandraUserRepository(session)

	taskService := services.NewTaskService(taskRepo)
	userService := services.NewUserService(userRepo)

	taskHandler := handlers.NewTaskHandler(taskService)
	welcomeHandler := handlers.NewWelcomeHandler()
	userHandler := handlers.NewUserHandler(userService, taskService)
	auth := middleware.NewAuthMiddleware(userService)

	mux := http.NewServeMux()

	mux.HandleFunc("/welcome", auth.Authenticate(welcomeHandler.Welcome))
	mux.HandleFunc("/printTasks", auth.Authenticate(taskHandler.GetAllTasksFromPjtHttp))
	mux.HandleFunc("/printProjects", auth.Authenticate(taskHandler.GetAllProjectsHttp))
	mux.HandleFunc("/writeTask", auth.Authenticate(taskHandler.WriteTaskHttp))
	mux.HandleFunc("/completeTask", auth.Authenticate(taskHandler.CompleteTaskHttp))
	mux.HandleFunc("/removeTask", auth.Authenticate(taskHandler.RemoveTaskHttp))
	mux.HandleFunc("/removeProject", auth.Authenticate(taskHandler.RemoveProjectHttp))
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
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
