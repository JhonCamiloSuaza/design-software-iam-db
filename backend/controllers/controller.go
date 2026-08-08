package controllers

import (
	"errors"
	"net/http"

	"iam-security-backend/httpx"
	"iam-security-backend/services"
)

func handleError(w http.ResponseWriter, err error) {
	var appErr *services.AppError
	if errors.As(err, &appErr) {
		httpx.Message(w, appErr.Status, appErr.Message)
		return
	}
	httpx.Message(w, http.StatusInternalServerError, "Ocurrio un error interno en el servidor.")
}
