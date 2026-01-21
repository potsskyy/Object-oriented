package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"todo/internal/handler"
	"todo/internal/middleware"
	"todo/internal/repository"
)

const (
	serverAddr = ":8080"
)

func main() {
	store := repository.NewMemoryRepo()

	userHandler := handler.NewAuthHandler(store)
	todoHandler := handler.NewTaskHandler(store)

	router := mux.NewRouter()

	router.HandleFunc("/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	todoRouter := router.PathPrefix("/").Subrouter()
	todoRouter.Use(middleware.WithAuth)

	todoRouter.HandleFunc("/add", todoHandler.Add).Methods("POST")
	todoRouter.HandleFunc("/update", todoHandler.Update).Methods("POST")
	todoRouter.HandleFunc("/resolve/{id}", todoHandler.Resolve).Methods("POST")
	todoRouter.HandleFunc("/delete/{id}", todoHandler.Delete).Methods("POST")
	todoRouter.HandleFunc("/get", todoHandler.GetAll).Methods("GET")
	todoRouter.HandleFunc("/get/{id}", todoHandler.GetByID).Methods("GET")
	todoRouter.HandleFunc("/archive", todoHandler.GetArchive).Methods("GET")

	log.Println("Server is running on", serverAddr)
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
