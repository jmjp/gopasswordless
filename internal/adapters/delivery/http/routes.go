package delivery

import (
	"hyperzoop/internal/adapters/delivery/http/controllers"
	"hyperzoop/internal/adapters/delivery/http/middlewares"
	repositories "hyperzoop/internal/adapters/repositories/pg"
	authusecases "hyperzoop/internal/usecases/auth"
)

func (s *HTTPServer) setupRoutes() {
	userRepository := repositories.NewUserPostgresRepository(s.db)
	magicRepository := repositories.NewMagicLinkPostgresRepository(s.db)
	sessionRepository := repositories.NewSessionPostgresRepository(s.db)

	loginUseCase := authusecases.NewLoginUseCase(userRepository, magicRepository)
	verifyUseCase := authusecases.NewVerifyUseCase(magicRepository, userRepository, sessionRepository)
	refreshUseCase := authusecases.NewRefreshUseCase(sessionRepository, userRepository)
	revokeUseCase := authusecases.NewRevokeUseCase(sessionRepository)
	sessionUseCase := authusecases.NewSessionUseCase(sessionRepository)

	authController := controllers.NewAuthenticationController(loginUseCase, verifyUseCase, refreshUseCase, revokeUseCase, sessionUseCase)

	s.router.Post("/auth/login", authController.Login)
	s.router.Get("/auth/verify", authController.Verify)
	s.router.Put("/auth/logout", middlewares.AutheMiddleware(authController.Logout))
	s.router.Post("/auth/refresh", authController.Refresh)
	s.router.Get("/auth/session", middlewares.AutheMiddleware(authController.Sessions))

}
