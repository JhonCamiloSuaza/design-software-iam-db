package controllers

import (
	"net/http"

	"iam-security-backend/httpx"
	"iam-security-backend/middleware"
	"iam-security-backend/models"
	"iam-security-backend/services"
)

type AuthController struct{ service *services.AuthService }

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}
func (controller *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterInput
	if !httpx.Decode(w, r, &input) {
		return
	}
	token, err := controller.service.Register(r.Context(), input)
	if err != nil {
		handleError(w, err)
		return
	}
	httpx.Write(w, http.StatusCreated, map[string]any{"message": "Usuario creado correctamente.", "token": token})
}
func (controller *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var input models.LoginInput
	if !httpx.Decode(w, r, &input) {
		return
	}
	token, err := controller.service.Login(r.Context(), input, r.RemoteAddr, r.UserAgent())
	if err != nil {
		handleError(w, err)
		return
	}
	httpx.Write(w, http.StatusOK, map[string]any{"message": "Inicio de sesion correcto.", "token": token})
}
func (controller *AuthController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var input models.ForgotPasswordInput
	if !httpx.Decode(w, r, &input) {
		return
	}
	token, err := controller.service.ForgotPassword(r.Context(), input.Email, r.RemoteAddr)
	if err != nil {
		handleError(w, err)
		return
	}
	response := map[string]string{"message": "Si el correo existe, se genero una solicitud de recuperacion."}
	if token != "" {
		response["message"] = "Solicitud creada. En este prototipo el token se muestra en pantalla; en produccion se enviaria por correo."
		response["resetToken"] = token
	}
	httpx.Write(w, http.StatusOK, response)
}
func (controller *AuthController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input models.ResetPasswordInput
	if !httpx.Decode(w, r, &input) {
		return
	}
	if err := controller.service.ResetPassword(r.Context(), input); err != nil {
		handleError(w, err)
		return
	}
	httpx.Message(w, http.StatusOK, "Contrasena actualizada. Ya puede iniciar sesion.")
}
func (controller *AuthController) Profile(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	profile, err := controller.service.Profile(r.Context(), claims.UserID)
	if err != nil {
		handleError(w, err)
		return
	}
	httpx.Write(w, http.StatusOK, profile)
}
func (controller *AuthController) Summary(w http.ResponseWriter, r *http.Request) {
	summary, err := controller.service.Summary(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	httpx.Write(w, http.StatusOK, summary)
}
