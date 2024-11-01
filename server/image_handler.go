package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var imageQueue = make(chan string, 100)

func init() {
	go processImageQueue()
}

// UploadHandler handles image uploads.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to retrieve image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "./images"
	os.MkdirAll(uploadDir, os.ModePerm)

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
	filePath := filepath.Join(uploadDir, fileName)
	saveFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer saveFile.Close()

	if _, err := io.Copy(saveFile, file); err != nil {
		http.Error(w, "Failed to copy image data", http.StatusInternalServerError)
		return
	}

	select {
	case imageQueue <- filePath:
		log.Printf("Image queued for analysis: %s", filePath)
	default:
		log.Println("Image queue is full, skipping analysis")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"status":"success","filePath":"%s"}`, filePath)))
}

func processImageQueue() {
	for filePath := range imageQueue {
		log.Printf("Processing image: %s", filePath)
		time.Sleep(2 * time.Second)
		log.Printf("Image processed: %s", filePath)
	}
}
