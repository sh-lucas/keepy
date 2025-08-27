package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// ANSI color codes for light colors
const (
	WHITE      = "\033[97m"       // Bright white
	PINK       = "\033[95m"       // Bright magenta/pink
	CYAN       = "\033[96m"       // Bright cyan
	PURPLE     = "\033[35m"       // Magenta/purple
	RED        = "\033[91m"       // Bright red
	YELLOW     = "\033[93m"       // Bright yellow
	ORANGE     = "\033[38;5;214m" // Light orange
	LAVENDER   = "\033[38;5;183m" // Lavender
	LIGHT_GREY = "\033[37m"       // Light grey
	RESET      = "\033[0m"        // Reset to default
)

// Global color palette
var colors = []string{WHITE, PINK, CYAN, PURPLE, RED, YELLOW, ORANGE, LAVENDER, LIGHT_GREY}

// LogEntry represents a log entry to be written
type LogEntry struct {
	Content string
}

// Channel for log entries
var logChannel chan LogEntry

// getColorByIP returns a color based on the IP address using hash
func getColorByIP(ip string) string {
	// Clean the IP address (remove port if present)
	cleanIP := strings.Split(ip, ":")[0]

	// Create hash of the IP
	h := fnv.New32a()
	h.Write([]byte(cleanIP))
	hash := h.Sum32()

	// Use hash to select color
	return colors[hash%uint32(len(colors))]
}

// fileWriter is a goroutine that handles all file writing operations
func fileWriter() {
	for entry := range logChannel {
		if OpenFile != nil {
			OpenFile.WriteString(entry.Content)
		} else {
			log.Println("Fuck the file closed for nothing")
		}
		// Also print to console
		fmt.Print(entry.Content)
	}
}

func main() {
	// Initialize the log channel
	logChannel = make(chan LogEntry, 1000) // Buffer of 1000 entries

	// Start the file writer goroutine
	go fileWriter()

	http.HandleFunc("/log", LogHandler)
	fmt.Println("Running super cool log server at port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}

var OpenFile *os.File

func init() {
	os.Mkdir("app", 0777)
	f, err := os.OpenFile("log_local.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Could not open the file wtf???")
	}
	OpenFile = f
}

func LogHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "You sent invalid stuff bro!", 400)
		return
	}

	// Get client IP address
	clientIP := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	}

	// Get color based on IP
	color := getColorByIP(clientIP)

	// Check if body content already ends with newline
	bodyStr := string(body)
	var newline string
	if len(bodyStr) > 0 && bodyStr[len(bodyStr)-1] != '\n' {
		newline = "\n"
	}

	// Create colored log entry content
	content := fmt.Sprintf("%s%s -> %s%s%s", color, time.Now().Format("2006/01/02 15:04:05"), bodyStr, RESET, newline)

	// Send to channel for async writing
	select {
	case logChannel <- LogEntry{Content: content}:
		// Successfully sent to channel
	default:
		// Channel is full, handle gracefully
		log.Println("Log channel is full, dropping log entry")
	}
}
