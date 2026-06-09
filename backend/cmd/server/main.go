package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "go_web/docs"
	"go_web/internal/bootstrap"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	casbinModelPath := flag.String("casbin-model", "configs/casbin_model.conf", "path to casbin model")
	flag.Parse()

	app, err := bootstrap.New(context.Background(), *configPath, *casbinModelPath)
	if err != nil {
		log.Fatalf("bootstrap application: %v", err)
	}
	defer app.Close()

	addr := fmt.Sprintf("%s:%d", app.Config.Server.Host, app.Config.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: app.Router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		app.Logger.Info(fmt.Sprintf("server starting on %s", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("run server: %v", err)
		}
	}()

	<-quit
	app.Logger.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		app.Logger.Error(fmt.Sprintf("server forced to shutdown: %v", err))
	}

	app.Logger.Info("server stopped")
}
