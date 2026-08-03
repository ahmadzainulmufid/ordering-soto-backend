package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"SotoAyam/config"
	"SotoAyam/internal/handlers"
	"SotoAyam/internal/repository"
	"SotoAyam/internal/routes"
	"SotoAyam/internal/services"
	"SotoAyam/internal/utils"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatalf(
			"gagal membaca environment: %v",
			err,
		)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf(
			"konfigurasi tidak valid: %v",
			err,
		)
	}

	ctx := context.Background()

	dbPool, err := config.NewPostgresPool(
		ctx,
		cfg.Database,
	)
	if err != nil {
		log.Fatalf(
			"koneksi PostgreSQL gagal: %v",
			err,
		)
	}
	defer dbPool.Close()

	jwtManager := utils.NewJWTManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.Issuer,
		cfg.JWT.AccessExpiration,
	)

	userRepository := repository.NewUserRepository(dbPool)

	authService := services.NewAuthService(
		userRepository,
		jwtManager,
	)

	authHandler := handlers.NewAuthHandler(
		authService,
	)

	router := routes.NewRouter(
		authHandler,
		jwtManager,
	)

	server := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf(
			"API berjalan pada http://localhost:%s",
			cfg.App.Port,
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf(
				"server gagal berjalan: %v",
				err,
			)
		}
	}()

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(
		signalChannel,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-signalChannel

	log.Println("mematikan server...")

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		log.Printf(
			"gagal mematikan server dengan aman: %v",
			err,
		)
	}

	log.Println("server berhenti")
}