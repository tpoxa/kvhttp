package http

import (
	"context"
	"log"
	"net/http"
)

func Start(ctx context.Context, listenAddr string, handler http.Handler) chan struct{} {
	done := make(chan struct{})
	// middleware
	// authentication
	srv := &http.Server{
		Addr:    listenAddr,
		Handler: handler,
	}

	go func() {
		log.Printf("listening: %s", listenAddr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("listenAndServe: " + err.Error())
		}
	}()

	go func() {
		defer close(done)

		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("shutdown: " + err.Error())
		}
	}()
	return done
}
