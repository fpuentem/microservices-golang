package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fpuentem/microservices-golang/ep4/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)

	// create the handlers
	ph := handlers.NewProducts(l)
	gh := handlers.NewGoodbye(l)

	// create a new serve mux and register for handlers
	sm := http.NewServeMux()
	sm.Handle("/", ph)
	sm.Handle("/goodbey", gh)

	// create a new server
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}

// https://www.youtube.com/watch?v=hodOppKJm5Y&t=564s
