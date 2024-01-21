package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/codebarz/go-micorservices/handlers"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	gp := handlers.NewProduct(l)

	ss := mux.NewRouter()

	getRouter := ss.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", gp.GetProducts)

	putRouter := ss.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", gp.UpdateProduct)
	putRouter.Use(gp.JSONValidationMiddleware)

	postRouter := ss.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", gp.AddProduct)
	postRouter.Use(gp.JSONValidationMiddleware)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      ss,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()

		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	l.Println("Recieved signal", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)
}
