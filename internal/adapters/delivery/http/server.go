package delivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"go.uber.org/zap"
)

type HTTPServer struct {
	port   int
	db     *sql.DB
	router *chi.Mux
}

func NewHTTPServer(port int, db *sql.DB) *HTTPServer {
	router := chi.NewRouter()
	return &HTTPServer{
		port:   port,
		db:     db,
		router: router,
	}
}
func (s *HTTPServer) Start() {
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(s.port),
		Handler: s.router,
	}
	s.setupMiddlewares()
	s.setupRoutes()

	go func() {
		fmt.Printf("Starting HTTP server on port %d\n", s.port)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("HTTP server error", zap.Error(err))
		}
		zap.L().Info("HTTP server stoping serving connections")
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	ctx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Fatal("HTTP server error", zap.Error(err))
	}
	zap.L().Info("HTTP server shutdown")

}

func (s *HTTPServer) setupMiddlewares() {
	s.router.Use(httprate.LimitByIP(100, 1*time.Minute))
	s.router.Use(middleware.CleanPath)
	if os.Getenv("env") == "dev" {
		s.router.Use(middleware.Logger)
		s.router.Use(cors.AllowAll().Handler)
	} else {
		s.router.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*.hyperzoop.com"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}
	s.router.Use(middleware.Recoverer)
}
