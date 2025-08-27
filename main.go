package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	cmdName := os.Args[1]
	logPath := os.Args[2]
	cmdArgs := os.Args[3:]

	lines := make(chan []byte, 5000) // buffer grande pra não travar leitura

	// Worker único pra processar linhas, linha a linha, e enviar pro servidor remoto
	go func() {
		for line := range lines {
			log.Println(string(line))
			// Envia para servidor remoto
			resp, err := http.Post("http://137.131.149.96/log/"+logPath, "text/plain", bytes.NewReader(line))
			if err != nil || resp.StatusCode != 200 {
				log.Printf("Ocorreu um erro enviando a request: %v, status: %d", err, resp.StatusCode)
			}
			// Adiciona delay de 100ms
			time.Sleep(100 * time.Millisecond)
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
		cmd.Wait()
		// se morrer, reinicia
		// depois de 5 segundos
		time.Sleep(5 * time.Second)
		log.Println("keepy reiniciando servidor...")
	}
}

func streamToChan(r io.Reader, ch chan<- []byte) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ch <- append([]byte{}, scanner.Bytes()...)
	}
}
