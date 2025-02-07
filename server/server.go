package main

import (
	_ "embed"
	"fmt"
	"foobar/database"
	"foobar/routes"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv" // Carrega .env
)

//go:generate templ generate

var (
	Port              = "8888"
	Addr              = "localhost"
	s    *http.Server = nil
)

func MustStartServer() {
	if s != nil {
		return
	}
	ready := make(chan bool, 1)
	go mustStartServer(ready)
	<-ready
}

func mustStartServer(ready chan<- bool) {

	if s != nil {
		ready <- true
		return
	}

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		Port = port
	}
	if addr := os.Getenv("SERVER_ADDR"); addr != "" {
		Addr = addr
	}

	_ = database.New()
	routes := routes.Register()

	s = &http.Server{
		Addr:           fmt.Sprintf("%s:%s", Addr, Port),
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Abre uma conexão de socket mas não bloqueia o fluxo do programa.
	// Isso significa que o servidor já está recebendo requests.
	// Isso é necessário para rodar os testes.
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}
	ready <- true
	fmt.Printf("Estamos ao vivo! teste em http://%s:%s\n", Addr, Port)

	log.Fatal(s.Serve(ln))
}

func main() {
	mustStartServer(make(chan bool, 1))
}
