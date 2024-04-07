package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"banner-service/internal/adapter/postgres"
	api "banner-service/internal/api/http"
	"banner-service/internal/config"
	"banner-service/internal/service"
)

func Start() error {
	log.Printf("start app")

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	db, err := initStorage(cfg.Database)
	if err != nil {
		return err
	}
	serv := service.NewService(db)
	handler := api.NewHandler(serv)

	httpServer := initHTTPServer(cfg.Server, handler)

	log.Printf("starting HTTP server on %s", httpServer.Addr)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("failed listen and serve: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server")
	if err := httpServer.Shutdown(context.Background()); err != nil {
		return err
	}
	if err := db.Close(context.Background()); err != nil {
		return err
	}

	return nil
}

func initStorage(cfg config.PostgresConfig) (service.BannerStorage, error) {
	log.Println("init postgres storage")
	return postgres.NewStorage(cfg)
}

func initHTTPServer(cfg config.HTTPConfig, handler *api.Handler) *http.Server {
	log.Println("init http server")
	return &http.Server{
		Addr:    fmt.Sprintf("[::]:%s", cfg.Port),
		Handler: api.InitRouter(handler),
	}
}
