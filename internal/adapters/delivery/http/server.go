package delivery

import (
	"context"
	"database/sql"
	"log"
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
	"github.com/redis/go-redis/v9"
)

type HTTPServer struct {
	port   int
	db     *sql.DB
	redis  *redis.Client
	router *chi.Mux
}

func NewHTTPServer(port int, db *sql.DB, redis *redis.Client) *HTTPServer {
	router := chi.NewRouter()
	return &HTTPServer{
		port:   port,
		db:     db,
		redis:  redis,
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
		log.Printf("Starting HTTP server on port %d\n", s.port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Printf("HTTP server stopping serving connections\n")
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Minute)
	defer shutdown()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
	log.Printf("HTTP server shutdown\n")
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
