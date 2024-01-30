package delivery

import (
	"hyperzoop/internal/adapters/delivery/http/controllers"
	"hyperzoop/internal/adapters/delivery/http/middlewares"
	repositories "hyperzoop/internal/adapters/repositories/pg"
	redisRepositories "hyperzoop/internal/adapters/repositories/redis"
	"hyperzoop/internal/core/ports"
	"hyperzoop/internal/core/services"
	"os"
)

func (s *HTTPServer) setupRoutes() {
	userRepository := repositories.NewUserPostgresRepository(s.db)
	var magicRepository ports.MagicLinkRepository
	if os.Getenv("env") == "prod" {
		magicRepository = redisRepositories.NewMagicLinkRedisRepository(s.redis)
	} else {
		magicRepository = repositories.NewMagicLinkPostgresRepository(s.db)
	}
	sessionRepository := repositories.NewSessionPostgresRepository(s.db)

	authService := services.NewAuthService(userRepository, magicRepository, sessionRepository)
	authController := controllers.NewAuthenticationController(authService)

	s.router.Post("/auth/login", authController.Login)
	s.router.Get("/auth/verify", authController.Verify)
	s.router.Put("/auth/logout", middlewares.AutheMiddleware(authController.Logout))
	s.router.Post("/auth/refresh", authController.Refresh)
	s.router.Get("/auth/session", middlewares.AutheMiddleware(authController.Sessions))

}
