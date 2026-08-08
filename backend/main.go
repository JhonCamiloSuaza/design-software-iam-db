package main

import (
	"context"
	"log"
	"net/http"

	"iam-security-backend/config"
	"iam-security-backend/controllers"
	"iam-security-backend/database"
	"iam-security-backend/middleware"
	"iam-security-backend/repositories"
	"iam-security-backend/routes"
	"iam-security-backend/services"
)

func main() {
	settings := config.Load()
	db, err := database.Connect(context.Background(), settings.DatabaseURL)
	if err != nil {
		log.Fatalf("No fue posible conectar con PostgreSQL: %v", err)
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepository, settings.JWTSecret)
	userService := services.NewUserService(userRepository)
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	systemController := controllers.NewSystemController(userRepository)
	authMiddleware := middleware.NewAuthMiddleware(settings.JWTSecret, userRepository)
	handler := routes.New(authController, userController, systemController, authMiddleware, settings.FrontendURL)

	log.Printf("IAM Security API listening on http://localhost:%s", settings.Port)
	log.Fatal(http.ListenAndServe(":"+settings.Port, handler))
}
