package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/razedwell/go-hand/internal/config"
	"github.com/razedwell/go-hand/internal/platform/cache"
	"github.com/razedwell/go-hand/internal/platform/db"
	"github.com/razedwell/go-hand/internal/platform/logger"
	"github.com/razedwell/go-hand/internal/postgres"
	"github.com/razedwell/go-hand/internal/security"
	authsrvc "github.com/razedwell/go-hand/internal/service/auth"
	"github.com/razedwell/go-hand/internal/service/user"
	transporthttp "github.com/razedwell/go-hand/internal/transport/http"
	"github.com/razedwell/go-hand/internal/transport/http/handler/auth"
	"github.com/razedwell/go-hand/internal/transport/http/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger.Init()
	cfg := config.LoadConfig()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbUrl := db.BuildDBUrl(cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode, cfg.Timezone)
	db, err := db.NewClient(dbUrl)
	if err != nil {
		logger.Log.Printf("Failed to connect to database: %s", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       int(cfg.RedisDB),
	})

	rdb := &cache.RedisClient{Client: redisClient}

	tokenRepo := postgres.NewTokenRepo(db)

	jwtManager := security.NewJWTManager(cfg.JWTAccessSecret, cfg.JWTRefreshSecret, time.Minute*time.Duration(cfg.JWTAccessExpiryMinutes), time.Hour*time.Duration(cfg.JWTRefreshExpiryHours), tokenRepo, rdb)
	authMW := middleware.Auth(jwtManager)

	userRepo := postgres.NewUserRepo(db)
	userService := user.NewService(userRepo)
	authService := authsrvc.NewService(userRepo, jwtManager)
	authHandler := auth.NewHandler(userService, authService, authMW)

	server := transporthttp.NewServer(":"+cfg.Port, authHandler)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Println("Server error:", err)
		}
	}()

	<-ctx.Done()
	logger.Log.Println("Graceful shutdown initiated...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutDown); err != nil {
		logger.Log.Fatalf("Server forced to close: %s", err)
		server.Close()
	} else {
		logger.Log.Println("Server exited properly.")
	}
	//Graceful shutdown logic can be added here if needed
}
