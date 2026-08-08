package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"iam-security-backend/httpx"
	"iam-security-backend/models"
	"iam-security-backend/repositories"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const claimsKey contextKey = "claims"

type AuthMiddleware struct {
	jwtSecret []byte
	repo      *repositories.UserRepository
}

func NewAuthMiddleware(jwtSecret string, repo *repositories.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: []byte(jwtSecret), repo: repo}
}
func ClaimsFromContext(ctx context.Context) models.Claims {
	claims, _ := ctx.Value(claimsKey).(models.Claims)
	return claims
}
func (middleware *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			httpx.Message(w, http.StatusUnauthorized, "Token de acceso requerido.")
			return
		}
		rawToken := strings.TrimPrefix(header, "Bearer ")
		claims := models.Claims{}
		token, err := jwt.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("metodo de firma invalido")
			}
			return middleware.jwtSecret, nil
		})
		if err != nil || !token.Valid {
			httpx.Message(w, http.StatusUnauthorized, "Token invalido o vencido.")
			return
		}
		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (middleware *AuthMiddleware) RequirePermission(featureCode string, next http.Handler) http.Handler {
	return middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := ClaimsFromContext(r.Context())
		if !middleware.repo.HasPermission(r.Context(), claims.UserID, featureCode) {
			httpx.Message(w, http.StatusForbidden, "No tiene permiso para administrar usuarios.")
			return
		}
		next.ServeHTTP(w, r)
	}))
}
