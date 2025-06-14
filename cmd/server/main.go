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
	"github.com/gorilla/mux"
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

	r := mux.NewRouter()

	r.HandleFunc("/welcome", auth.Authenticate(welcomeHandler.Welcome)).Methods("GET", "OPTIONS")
	r.HandleFunc("/printTasks", auth.Authenticate(taskHandler.GetAllTasksFromPjtHttp)).Methods("GET", "OPTIONS")
	r.HandleFunc("/printProjects", auth.Authenticate(taskHandler.GetAllProjectsHttp)).Methods("GET", "OPTIONS")
	r.HandleFunc("/writeTask", auth.Authenticate(taskHandler.WriteTaskHttp)).Methods("POST", "OPTIONS")
	r.HandleFunc("/completeTask", auth.Authenticate(taskHandler.CompleteTaskHttp)).Methods("GET", "OPTIONS")
	r.HandleFunc("/removeTask", auth.Authenticate(taskHandler.RemoveTaskHttp)).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/removeProject", auth.Authenticate(taskHandler.RemoveProjectHttp)).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/deactivate", auth.Authenticate(userHandler.DeleteUser)).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/register", userHandler.Register).Methods("POST", "OPTIONS")

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "7071" // Default port
	}
	serverAddr := ":" + serverPort
	handler := middleware.CORS(r)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: handler,
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
