package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func initialRouteHandler(responseWritter http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	responseWritter.Write([]byte(createGreeting(query.Get("name"))))
}

func createGreeting(name string) string {
	if name == "" {
		name = "Guest"
	}
	return "Hello, it works - " + name + " - And this is a new application version\n"
}

func main() {
	muxRouterHandler := mux.NewRouter()

	muxRouterHandler.HandleFunc("/", initialRouteHandler)

	muxServer := buildMuxServer(muxRouterHandler)

	startServer(muxServer)

	waitForGracefulShutDown(muxServer)
}

func buildMuxServer(muxRouterHandler http.Handler) *http.Server {
	httpServer := &http.Server{
		Handler:      muxRouterHandler,
		Addr:         ":3000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return httpServer
}
func startServer(muxServer *http.Server) {
	log.Println("Starting Gorilla Server")
	if err := muxServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func waitForGracefulShutDown(muxServer *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	muxServer.Shutdown(ctx)

	os.Exit(0)
}
