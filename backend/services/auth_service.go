package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"iam-security-backend/models"
	"iam-security-backend/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repositories.UserRepository
	jwtSecret []byte
}

func NewAuthService(repo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: []byte(jwtSecret)}
}
func normalizeEmail(email string) string { return strings.ToLower(strings.TrimSpace(email)) }
func tokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (service *AuthService) createJWT(userID, email string) (string, error) {
	claims := models.Claims{UserID: userID, Email: email, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), IssuedAt: jwt.NewNumericDate(time.Now())}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(service.jwtSecret)
}

func (service *AuthService) Register(ctx context.Context, input models.RegisterInput) (string, error) {
	input.Email = normalizeEmail(input.Email)
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	if input.Email == "" || input.FirstName == "" || input.LastName == "" || len(input.Password) < 8 {
		return "", NewError(http.StatusBadRequest, "Complete los datos y use una contrasena de minimo 8 caracteres.")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	userID, err := service.repo.Register(ctx, input, string(hash))
	if errors.Is(err, repositories.ErrDuplicateEmail) {
		return "", NewError(http.StatusConflict, "Ya existe una cuenta con este correo.")
	}
	if err != nil {
		return "", err
	}
	return service.createJWT(userID, input.Email)
}

func (service *AuthService) Login(ctx context.Context, input models.LoginInput, ip, userAgent string) (string, error) {
	input.Email = normalizeEmail(input.Email)
	user, err := service.repo.FindByEmail(ctx, input.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		service.repo.RecordLogin(ctx, nil, input.Email, "USER_NOT_FOUND", ip, userAgent)
		return "", NewError(http.StatusUnauthorized, "Correo o contrasena incorrectos.")
	}
	if err != nil {
		return "", err
	}
	if !user.IsActive {
		service.repo.RecordLogin(ctx, &user.ID, input.Email, "ACCOUNT_LOCKED", ip, userAgent)
		return "", NewError(http.StatusUnauthorized, "La cuenta se encuentra inactiva.")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		service.repo.RecordLogin(ctx, &user.ID, input.Email, "INVALID_PASSWORD", ip, userAgent)
		return "", NewError(http.StatusUnauthorized, "Correo o contrasena incorrectos.")
	}
	token, err := service.createJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}
	if err = service.repo.CompleteLogin(ctx, user, tokenHash(token), userAgent); err != nil {
		return "", err
	}
	service.repo.RecordLogin(ctx, &user.ID, user.Email, "SUCCESS", ip, userAgent)
	return token, nil
}

func randomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
func (service *AuthService) ForgotPassword(ctx context.Context, email, ip string) (string, error) {
	token, err := randomToken()
	if err != nil {
		return "", err
	}
	exists, err := service.repo.CreateResetRequest(ctx, normalizeEmail(email), tokenHash(token), ip)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", nil
	}
	return token, nil
}
func (service *AuthService) ResetPassword(ctx context.Context, input models.ResetPasswordInput) error {
	if input.Token == "" || len(input.Password) < 8 {
		return NewError(http.StatusBadRequest, "Token y contrasena de minimo 8 caracteres son obligatorios.")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = service.repo.ResetPassword(ctx, tokenHash(input.Token), string(hash))
	if errors.Is(err, pgx.ErrNoRows) {
		return NewError(http.StatusBadRequest, "El token no es valido o ya vencio.")
	}
	return err
}
func (service *AuthService) Profile(ctx context.Context, userID string) (models.Profile, error) {
	return service.repo.Profile(ctx, userID)
}
func (service *AuthService) Summary(ctx context.Context) (models.Summary, error) {
	return service.repo.Summary(ctx)
}
