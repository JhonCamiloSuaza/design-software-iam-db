package routes

import (
	"net/http"

	"iam-security-backend/controllers"
	"iam-security-backend/middleware"
)

func New(authController *controllers.AuthController, userController *controllers.UserController, systemController *controllers.SystemController, authMiddleware *middleware.AuthMiddleware, frontendURL string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", systemController.Health)
	mux.HandleFunc("POST /api/auth/register", authController.Register)
	mux.HandleFunc("POST /api/auth/login", authController.Login)
	mux.HandleFunc("POST /api/auth/forgot-password", authController.ForgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", authController.ResetPassword)
	mux.Handle("GET /api/me", authMiddleware.RequireAuth(http.HandlerFunc(authController.Profile)))
	mux.Handle("GET /api/catalog/summary", authMiddleware.RequireAuth(http.HandlerFunc(authController.Summary)))
	mux.Handle("GET /api/users", authMiddleware.RequirePermission("IDENTITY_USER_MANAGE", http.HandlerFunc(userController.List)))
	mux.Handle("POST /api/users", authMiddleware.RequirePermission("IDENTITY_USER_MANAGE", http.HandlerFunc(userController.Create)))
	mux.Handle("PUT /api/users/{id}", authMiddleware.RequirePermission("IDENTITY_USER_MANAGE", http.HandlerFunc(userController.Update)))
	mux.Handle("DELETE /api/users/{id}", authMiddleware.RequirePermission("IDENTITY_USER_MANAGE", http.HandlerFunc(userController.Deactivate)))
	mux.Handle("GET /api/roles", authMiddleware.RequirePermission("IDENTITY_USER_MANAGE", http.HandlerFunc(userController.Roles)))
	return middleware.CORS(frontendURL, mux)
}
