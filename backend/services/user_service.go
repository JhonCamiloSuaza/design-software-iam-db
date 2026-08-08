package services

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"iam-security-backend/models"
	"iam-security-backend/repositories"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{ repo *repositories.UserRepository }

func NewUserService(repo *repositories.UserRepository) *UserService { return &UserService{repo: repo} }
func validActorType(value string) bool {
	return value == "USER" || value == "INSTRUCTOR" || value == "LEARNER"
}
func cleanManagedInput(input models.ManagedUserInput) models.ManagedUserInput {
	input.Email = normalizeEmail(input.Email)
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.ActorType = strings.TrimSpace(input.ActorType)
	input.RoleName = strings.TrimSpace(input.RoleName)
	return input
}
func (service *UserService) List(ctx context.Context) ([]models.ManagedUser, error) {
	return service.repo.ListUsers(ctx)
}
func (service *UserService) Roles(ctx context.Context) ([]models.Role, error) {
	return service.repo.ListRoles(ctx)
}
func (service *UserService) Create(ctx context.Context, input models.ManagedUserInput, assignedBy string) error {
	input = cleanManagedInput(input)
	if input.Email == "" || input.FirstName == "" || input.LastName == "" || len(input.Password) < 8 || !validActorType(input.ActorType) || input.RoleName == "" {
		return NewError(http.StatusBadRequest, "Complete todos los campos; la contrasena debe tener 8 caracteres.")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = service.repo.CreateManaged(ctx, input, string(hash), assignedBy)
	if errors.Is(err, repositories.ErrDuplicateEmail) {
		return NewError(http.StatusConflict, "Ya existe un usuario con ese correo.")
	}
	if errors.Is(err, repositories.ErrRoleNotFound) {
		return NewError(http.StatusBadRequest, "El rol seleccionado no existe.")
	}
	return err
}
func (service *UserService) Update(ctx context.Context, userID string, input models.ManagedUserInput) error {
	input = cleanManagedInput(input)
	if input.FirstName == "" || input.LastName == "" || !validActorType(input.ActorType) || input.IsActive == nil {
		return NewError(http.StatusBadRequest, "Nombres, apellidos, tipo y estado son obligatorios.")
	}
	err := service.repo.UpdateManaged(ctx, userID, input)
	if errors.Is(err, pgx.ErrNoRows) {
		return NewError(http.StatusNotFound, "Usuario no encontrado.")
	}
	return err
}
func (service *UserService) Deactivate(ctx context.Context, userID, currentUserID string) error {
	if userID == currentUserID {
		return NewError(http.StatusBadRequest, "No puede desactivar su propia cuenta.")
	}
	err := service.repo.Deactivate(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return NewError(http.StatusNotFound, "Usuario no encontrado.")
	}
	return err
}
