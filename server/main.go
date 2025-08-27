package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/log/", LogHandler)
	http.HandleFunc("/log", LogHandler) // Fallback sem trailing slash
	fmt.Println("Opening server at port 80")
	err := http.ListenAndServe(":80", nil)
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

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "You sent invalid stuff bro!", 400)
		return
	}

	// Open file for appending (create if it doesn't exist)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", filename, err)
		http.Error(w, "Could not open file", 500)
		return
	}
	defer file.Close()

	bodyStr := strings.TrimRight(string(body), "\n")
	content := fmt.Sprintf("%s @ %s ---\n		%s\n", time.Now().Format(time.RFC3339), r.RemoteAddr, bodyStr)

	// Append body content to file
	_, err = fmt.Fprint(file, content)
	if err != nil {
		log.Printf("Error writing to file %s: %v\n", filename, err)
		http.Error(w, "Could not write to file", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Content appended to %s successfully", filename)
}
