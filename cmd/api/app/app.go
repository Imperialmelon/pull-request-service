package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Imperialmelon/AvitoTechTest/internal/app/handlers"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/repository"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/service"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/usecase"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Run() error {
	config := LoadConfig()

	store, err := repository.NewPostgresStore(config.GetDatabaseDSN())
	if err != nil {
		return err
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Fatalln("Failed to close database connection", "error", err)
		}
	}()

	var repo service.Repository = store
	serviceInstance := service.NewService(repo)
	usecaseInstance := usecase.NewUseCase(serviceInstance)
	handler := handlers.NewHandler(usecaseInstance)

	r := mux.NewRouter()
	handler.Register(r)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:    config.GetServerAddress(),
		Handler: r,
	}

	go func() {
		log.Println("Server running at", "address", config.GetServerAddress())
		log.Println("Starting HTTP server")
		err = srv.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Println("Server exited")
	return nil
}
