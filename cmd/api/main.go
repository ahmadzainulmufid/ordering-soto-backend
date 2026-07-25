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

	"my-backend/config"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("gagal membaca environment variable: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("konfigurasi aplikasi tidak valid: %v", err)
	}

	appCtx := context.Background()

	dbPool, err := config.NewPostgresPool(
		appCtx,
		cfg.Database,
	)
	if err != nil {
		log.Fatalf("koneksi database gagal: %v", err)
	}
	defer dbPool.Close()

	/*
		Repository nantinya menerima dbPool, misalnya:

		userRepository := repository.NewUserRepository(dbPool)
		productRepository := repository.NewProductRepository(dbPool)
		orderRepository := repository.NewOrderRepository(dbPool)
	*/

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		ctx, cancel := context.WithTimeout(
			r.Context(),
			3*time.Second,
		)
		defer cancel()

		if err := dbPool.Ping(ctx); err != nil {
			http.Error(
				w,
				`{"status":"error","database":"disconnected"}`,
				http.StatusServiceUnavailable,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(
			[]byte(`{"status":"ok","database":"connected"}`),
		)
	})

	server := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf(
			"%s berjalan pada http://localhost:%s",
			cfg.App.Name,
			cfg.App.Port,
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server gagal dijalankan: %v", err)
		}
	}()

	shutdownSignal := make(chan os.Signal, 1)

	signal.Notify(
		shutdownSignal,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-shutdownSignal

	log.Println("mematikan server...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server gagal dimatikan dengan aman: %v", err)
	}

	log.Println("server berhasil dimatikan")
}