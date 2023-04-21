package main

import (
	"context"
	"kvhttp/internal/hasher"
	"kvhttp/internal/http"
	"kvhttp/internal/router"
	"kvhttp/internal/store/mem"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	h := hasher.NewSha256()
	s := mem.NewStore(h)
	handler := router.NewRouter(s, h)
	done := http.Start(ctx, "0.0.0.0:8085", handler)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()
	log.Println("graceful shutdown")
	<-done
}
