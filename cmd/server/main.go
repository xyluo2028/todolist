package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "cassandra" // Default to cassandra
	}
	log.Printf("Using storage type: %s", storageType)

	var taskRepo repository.TaskRepository
	var userRepo repository.UserRepository

	if storageType == "cassandra" {
		cassandraHostsEnv := os.Getenv("CASSANDRA_HOSTS")
		if cassandraHostsEnv == "" {
			cassandraHostsEnv = "127.0.0.1:9042" // Default for local development
		}
		cassandraKeyspaceEnv := os.Getenv("CASSANDRA_KEYSPACE")
		if cassandraKeyspaceEnv == "" {
			cassandraKeyspaceEnv = "todolist" // Default keyspace
		}

		cluster := gocql.NewCluster(strings.Split(cassandraHostsEnv, ",")...) // Use configured hosts
		cluster.Keyspace = cassandraKeyspaceEnv                               // Use configured keyspace
		cluster.Consistency = gocql.Quorum
		cluster.Timeout = 5 * time.Second

		session, err := cluster.CreateSession()
		if err != nil {
			log.Fatalf("Could not connect to Cassandra: %v", err)
		}
		defer session.Close()
		log.Println("Successfully connected to Cassandra.")
		taskRepo = repository.NewCassandraTaskRepository(session)
		userRepo = repository.NewCassandraUserRepository(session)
	} else if storageType == "inmem" {
		log.Println("Using in-memory storage.")
		taskRepo = repository.NewInMemTaskRepository()
		userRepo = repository.NewInMemUserRepository()
	} else {
		log.Fatalf("Invalid STORAGE_TYPE: %s. Supported values are 'cassandra' or 'inmem'.", storageType)
	}

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

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "7071" // Default port
	}
	serverAddr := ":" + serverPort

	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		log.Printf("shutting down server..")
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	log.Printf("Server starting on %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}
