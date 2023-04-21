// Created with Strapit
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	HOST = os.Getenv("HOST")
	PORT = os.Getenv("PORT")
	ADDR = HOST + ":" + PORT
)

func main() {
	parent := context.Background()

	signals := []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	signaller := make(chan os.Signal, len(signals))
	signal.Notify(signaller, signals...)

	server := &http.Server{
		Addr:        ADDR,
		Handler:     http.FileServer(http.Dir("static")),
		BaseContext: func(l net.Listener) context.Context { return parent },
	}

	go func() {
		log.Printf("Serving HTTP from %s", ADDR)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	sig := <-signaller
	log.Printf("%s signal received, initiating graceful shutdown", sig.String())
	shutCtx, cancel := context.WithTimeout(parent, time.Second*5)
	defer cancel()
	err := server.Shutdown(shutCtx)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("HTTP(S) server shutdown complete")
}
