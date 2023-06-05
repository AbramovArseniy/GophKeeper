package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AbramovArseniy/GophKeeper/internal/server/handlers"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/config"
)

func StartServer() {
	cfg := config.SetServerParams()
	var err error
	if cfg.DatabaseAddress != "" {
		cfg.Database, err = sql.Open("pgx", cfg.DatabaseAddress)
		if err != nil {
			log.Println("opening DB error:", err)
			cfg.Database = nil
		}
		defer cfg.Database.Close()
	} else {
		cfg.Database = nil
	}
	s := handlers.NewServer(cfg)
	handler := s.Route()
	srv := &http.Server{
		Addr:    s.Addr,
		Handler: handler,
	}
	log.Printf("HTTP server started at %s", s.Addr)
	idleConnsClosed := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigs
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func main() {}
