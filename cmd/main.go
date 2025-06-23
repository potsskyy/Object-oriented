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
	runAddr = ":8080"
)

func main() {
	repo := repository.NewMemoryRepo()

	authHandler := handler.NewAuthHandler(repo)
	taskHandler := handler.NewTaskHandler(repo)

	r := mux.NewRouter()

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	taskRouter := r.PathPrefix("/").Subrouter()
	taskRouter.Use(middleware.WithAuth)

	taskRouter.HandleFunc("/add", taskHandler.Add).Methods("POST")
	taskRouter.HandleFunc("/update", taskHandler.Update).Methods("POST")
	taskRouter.HandleFunc("/resolve/{id}", taskHandler.Resolve).Methods("POST")
	taskRouter.HandleFunc("/delete/{id}", taskHandler.Delete).Methods("POST")
	taskRouter.HandleFunc("/get", taskHandler.GetAll).Methods("GET")
	taskRouter.HandleFunc("/get/{id}", taskHandler.GetByID).Methods("GET")
	taskRouter.HandleFunc("/archive", taskHandler.GetArchive).Methods("GET")

	log.Println("Server is running on", runAddr)
	if err := http.ListenAndServe(runAddr, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
