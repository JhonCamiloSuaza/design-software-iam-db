package controllers

import (
	"net/http"

	"iam-security-backend/httpx"
	"iam-security-backend/middleware"
	"iam-security-backend/models"
	"iam-security-backend/services"
)

type UserController struct{ service *services.UserService }

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}
func (controller *UserController) List(w http.ResponseWriter, r *http.Request) {
	users, err := controller.service.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	httpx.Write(w, http.StatusOK, users)
}
func (controller *UserController) Roles(w http.ResponseWriter, r *http.Request) {
	roles, err := controller.service.Roles(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	httpx.Write(w, http.StatusOK, roles)
}
func (controller *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var input models.ManagedUserInput
	if !httpx.Decode(w, r, &input) {
		return
	}
	claims := middleware.ClaimsFromContext(r.Context())
	if err := controller.service.Create(r.Context(), input, claims.UserID); err != nil {
		handleError(w, err)
		return
	}
	httpx.Message(w, http.StatusCreated, "Usuario creado y rol asignado.")
}
func (controller *UserController) Update(w http.ResponseWriter, r *http.Request) {
	var input models.ManagedUserInput
	if !httpx.Decode(w, r, &input) {
		return
	}
	if err := controller.service.Update(r.Context(), r.PathValue("id"), input); err != nil {
		handleError(w, err)
		return
	}
	httpx.Message(w, http.StatusOK, "Usuario actualizado correctamente.")
}
func (controller *UserController) Deactivate(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if err := controller.service.Deactivate(r.Context(), r.PathValue("id"), claims.UserID); err != nil {
		handleError(w, err)
		return
	}
	httpx.Message(w, http.StatusOK, "Usuario desactivado y sesiones revocadas.")
}
