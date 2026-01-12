package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/razedwell/go-hand/internal/model"
	"github.com/razedwell/go-hand/internal/platform/logger"
	"github.com/razedwell/go-hand/internal/service/user"
	"github.com/razedwell/go-hand/internal/transport/http/handler/auth"
	"github.com/razedwell/go-hand/internal/transport/http/middleware"
	"golang.org/x/crypto/bcrypt"
)

var mockHash, _ = bcrypt.GenerateFromPassword([]byte("<PASSWORD>"), bcrypt.DefaultCost)

// mockUserRepo implements the repository interface expected by the service
type mockUserRepo struct{}

func (m *mockUserRepo) FindUserByEmail(email string) (*model.User, error) {
	if email == "<EMAIL>" {
		return &model.User{
			ID:           1,
			Email:        "<EMAIL>",
			PasswordHash: string(mockHash),
			FirstName:    "John",
			LastName:     "Doe",
		}, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) FindUserById(id int64) (*model.User, error) {
	return nil, errors.New("user not found")
}

func main() {
	logger.Init()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	userService := user.NewService(&mockUserRepo{})
	authHandler := auth.NewHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("/", authHandler.Home)
	mux.HandleFunc("/login", authHandler.Login)
	handler := middleware.LoggingMiddleware(mux)
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Server error:", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
	server.Shutdown(ctx)
	//Graceful shutdown logic can be added here if needed
}
