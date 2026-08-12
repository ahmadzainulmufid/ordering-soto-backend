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
	// Load Environment

	if err := config.LoadEnv(); err != nil {
		log.Fatalf(
			"gagal membaca environment: %v",
			err,
		)
	}

	// Load Configuration

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf(
			"konfigurasi tidak valid: %v",
			err,
		)
	}

	// Context

	ctx := context.Background()

	// PostgreSQL Connection

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

	log.Println("koneksi PostgreSQL berhasil")

	// JWT Manager

	jwtManager := utils.NewJWTManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.Issuer,
		cfg.JWT.AccessExpiration,
	)

	// Repository

	userRepository :=
		repository.NewUserRepository(
			dbPool,
		)

	categoryRepository :=
		repository.NewCategoryRepository(
			dbPool,
		)

	productRepository :=
		repository.NewProductRepository(
			dbPool,
		)

	diningTableRepository :=
		repository.NewDiningTableRepository(
			dbPool,
		)

	orderRepository :=
		repository.NewOrderRepository(
			dbPool,
		)

	// Service

	authService :=
		services.NewAuthService(
			userRepository,
			jwtManager,
		)

	categoryService :=
		services.NewCategoryService(
			categoryRepository,
		)

	productService :=
		services.NewProductService(
			productRepository,
			categoryRepository,
		)

	diningTableService :=
		services.NewDiningTableService(
			diningTableRepository,
		)

	orderService :=
		services.NewOrderService(
			orderRepository,
			productRepository,
			diningTableRepository,
		)

	// Handler

	authHandler :=
		handlers.NewAuthHandler(
			authService,
		)

	categoryHandler :=
		handlers.NewCategoryHandler(
			categoryService,
		)

	productHandler :=
		handlers.NewProductHandler(
			productService,
		)

	diningTableHandler :=
		handlers.NewDiningTableHandler(
			diningTableService,
		)

	orderHandler :=
		handlers.NewOrderHandler(
			orderService,
		)

	// Router

	router := routes.NewRouter(
		routes.Handlers{
			AuthHandler:        authHandler,
			CategoryHandler:    categoryHandler,
			ProductHandler:     productHandler,
			DiningTableHandler: diningTableHandler,
			OrderHandler:       orderHandler,
		},
		jwtManager,
	)

	// HTTP Server

	server := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Run Server

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

	// Graceful Shutdown

	signalChannel := make(
		chan os.Signal,
		1,
	)

	signal.Notify(
		signalChannel,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-signalChannel

	log.Println(
		"mematikan server...",
	)

	shutdownContext, cancel :=
		context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

	defer cancel()

	if err := server.Shutdown(
		shutdownContext,
	); err != nil {

		log.Printf(
			"gagal mematikan server dengan aman: %v",
			err,
		)
	}

	log.Println(
		"server berhenti",
	)
}