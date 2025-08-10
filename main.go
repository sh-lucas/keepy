package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	lines := make(chan []byte, 1000) // buffer grande pra não travar leitura

	// Worker único pra processar linhas
	go func() {

		for line := range lines {
			log.Println(string(line))
			// Envia para servidor remoto
			resp, err := http.Post("http://137.131.149.96/log", "text/plain", bytes.NewReader(line))
			if err != nil || resp.StatusCode != 200 {
				log.Printf("Ocorreu um erro enviando a request: %v, status: %d", err, resp.StatusCode)
			}
		}
	}()

	for {
		cmd := exec.Command(cmdName, cmdArgs...)

		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()

		go streamToChan(stdout, lines)
		go streamToChan(stderr, lines)

		if err := cmd.Start(); err != nil {
			panic(err)
		}

		cmd.Start()
		cmd.Wait() // se morrer, reinicia
	}
}

func streamToChan(r io.Reader, ch chan<- []byte) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ch <- append([]byte{}, scanner.Bytes()...)
	}
}
