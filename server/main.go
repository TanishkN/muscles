package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"musclesfyi/server/models"

	"github.com/joho/godotenv"
)

func initializeDB() {
	if err := models.InitDB(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	fmt.Println("Connected to the PostgreSQL database successfully.")
}

func setupRoutes() {
	http.HandleFunc("/upload", UploadHandler)     // Endpoint for uploading images
	http.HandleFunc("/datafeed", GetPostsHandler) // Endpoint for viewing datafeed
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
}

var cache = NewCache() // Initialize in-memory cache for recent posts

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetRecentPosts(100) // Fetch from database
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	initializeDB()
	defer models.DB.Close()

	setupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
