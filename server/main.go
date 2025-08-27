package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Setup logging to file
	logDir := "/app"
	logFilename := "server_log.txt"
	logPath := filepath.Join(logDir, logFilename)

	// Ensure the log directory exists
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		log.Printf("Failed to create log directory %s: %v", logDir, err)
	} else {
		// Open log file for appending (create if it doesn't exist)
		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Printf("Failed to open log file %s: %v", logPath, err)
		} else {
			// Set log output to both stdout and file
			multiWriter := io.MultiWriter(os.Stdout, logFile)
			log.SetOutput(multiWriter)
			log.Printf("Server logging initialized. Logs will be written to %s", logPath)
		}
	}

	http.HandleFunc("/log/", LogHandler)
	http.HandleFunc("/log", LogHandler) // Fallback sem trailing slash
	fmt.Println("Opening server at port 80")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}

func LogHandler(w http.ResponseWriter, r *http.Request) {
	var filename string

	// Extract source parameter from URL path
	path := strings.TrimPrefix(r.URL.Path, "/log")
	path = strings.TrimPrefix(path, "/") // Remove leading slash if present

	if path == "" || path == "/" {
		// No parameter provided, use fallback
		filename = "log_local.txt"
		log.Println("No source parameter provided, using fallback: log_local.txt")
	} else {
		// Use the provided parameter as filename
		// Remove any extra slashes and get just the first part
		pathParts := strings.Split(path, "/")
		filename = pathParts[0]
		log.Printf("Using source parameter: %s\n", filename)
	}

	// Create full path to log file in /app directory
	logDir := "/app"
	filePath := filepath.Join(logDir, filename)

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "You sent invalid stuff bro!", 400)
		return
	}

	// Ensure the log directory exists
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		log.Printf("Error creating directory %s: %v\n", logDir, err)
		http.Error(w, "Could not create log directory", 500)
		return
	}

	// Open file for appending (create if it doesn't exist)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", filePath, err)
		http.Error(w, "Could not open file", 500)
		return
	}
	defer file.Close()

	bodyStr := strings.TrimRight(string(body), "\n")
	content := fmt.Sprintf("%s @ %s --> %s\n", time.Now().Format(time.RFC3339), r.RemoteAddr, bodyStr)

	// Append body content to file
	_, err = fmt.Fprint(file, content)
	if err != nil {
		log.Printf("Error writing to file %s: %v\n", filePath, err)
		http.Error(w, "Could not write to file", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Content appended to %s successfully", filePath)
}
