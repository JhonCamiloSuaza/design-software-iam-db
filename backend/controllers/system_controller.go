package controllers

import (
	"net/http"

	"iam-security-backend/httpx"
	"iam-security-backend/repositories"
)

type SystemController struct{ repo *repositories.UserRepository }

func NewSystemController(repo *repositories.UserRepository) *SystemController {
	return &SystemController{repo: repo}
}
func (controller *SystemController) Health(w http.ResponseWriter, r *http.Request) {
	if err := controller.repo.Ping(r.Context()); err != nil {
		httpx.Write(w, http.StatusServiceUnavailable, map[string]string{"status": "database unavailable"})
		return
	}
	httpx.Write(w, http.StatusOK, map[string]string{"status": "ok", "service": "iam-security-go"})
}
