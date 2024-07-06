package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// fmt.Println("Starting Main App")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := chi.NewRouter()
	// router.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello World"))
	// })
	router.Mount("/task", GetAllRoutes())

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("server running on port:", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal("error while starting server:", err)
	}
	// log.Println("server running on port: 8100")
}
