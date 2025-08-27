package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/log", LogHandler)
	fmt.Println("Opening server at port 80")
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
	}

	log.Println(string(body) + "\n")

	if OpenFile != nil {
		OpenFile.Write(append(body, '\n'))
	} else {
		log.Println("Fuck the file closed for nothing")
	}
}
