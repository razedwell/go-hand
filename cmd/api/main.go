package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/razedwell/go-hand/internal/config"
	"github.com/razedwell/go-hand/internal/platform/logger"
	"github.com/razedwell/go-hand/internal/platform/postgres"
	pgrepo "github.com/razedwell/go-hand/internal/postgres"
	"github.com/razedwell/go-hand/internal/security"
	authsrvc "github.com/razedwell/go-hand/internal/service/auth"
	"github.com/razedwell/go-hand/internal/service/user"
	transporthttp "github.com/razedwell/go-hand/internal/transport/http"
	"github.com/razedwell/go-hand/internal/transport/http/handler/auth"
	"github.com/razedwell/go-hand/internal/transport/http/middleware"
)

// var mockHash, _ = bcrypt.GenerateFromPassword([]byte("<PASSWORD>"), bcrypt.DefaultCost)

// // mockUserRepo implements the repository interface expected by the service
// type mockUserRepo struct{}

// func (m *mockUserRepo) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
// 	if email == "<EMAIL>" {
// 		return &model.User{
// 			ID:           1,
// 			Email:        "<EMAIL>",
// 			PasswordHash: string(mockHash),
// 			FirstName:    "John",
// 			LastName:     "Doe",
// 		}, nil
// 	}
// 	return nil, errors.New("user not found")
// }

// func (m *mockUserRepo) FindUserById(ctx context.Context, id int64) (*model.User, error) {
// 	return nil, errors.New("user not found")
// }

func main() {
	logger.Init()
	cfg := config.LoadConfig()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbUrl := postgres.BuildDBUrl(cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)
	db, err := postgres.NewClient(dbUrl)
	if err != nil {
		logger.Log.Printf("Failed to connect to database: %s", err)
	}

	jwtManager := security.NewJWTManager(cfg.JWTSecret, time.Minute*15)
	authMW := middleware.Auth(jwtManager)

	userRepo := pgrepo.NewUserRepo(db)
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
