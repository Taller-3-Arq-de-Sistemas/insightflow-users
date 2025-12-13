package auth

import (
	"net/http"

	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/api"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/json"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/validator"
)

type Handler struct {
	service *svc
}

func NewHandler(service *svc) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var params LoginParams
	if err := json.Read(r, &params); err != nil {
		api.BadRequest(w, err.Error())
		return
	}

	if errors := validator.ValidateStruct(params); errors != nil {
		api.ValidationErrors(w, errors)
		return
	}

	token, err := h.service.Login(r.Context(), params)
	if err != nil {
		if err == ErrInvalidCredentials {
			api.InvalidCredentials(w, err.Error())
			return
		}
		api.InternalServerError(w, err)
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var params RegisterParams
	if err := json.Read(r, &params); err != nil {
		api.BadRequest(w, err.Error())
		return
	}

	if errors := validator.ValidateStruct(params); errors != nil {
		api.ValidationErrors(w, errors)
		return
	}

	token, err := h.service.Register(r.Context(), params)
	if err != nil {
		if err == ErrInvalidCredentials || err == ErrInvalidBirthDate {
			api.BadRequest(w, err.Error())
			return
		}
		if err == ErrUserAlreadyExists {
			api.Conflict(w, err.Error())
			return
		}
		api.InternalServerError(w, err)
		return
	}

	json.Write(w, http.StatusCreated, map[string]string{"token": token})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Logout(r.Context()); err != nil {
		api.InternalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var params ValidateTokenParams
	if err := json.Read(r, &params); err != nil {
		api.BadRequest(w, err.Error())
		return
	}

	if errors := validator.ValidateStruct(params); errors != nil {
		api.ValidationErrors(w, errors)
		return
	}

	data, err := h.service.ValidateToken(r.Context(), params)
	if err != nil {
		if err == ErrInvalidToken {
			api.Unauthorized(w, err.Error())
			return
		}
		api.InternalServerError(w, err)
		return
	}
	json.Write(w, http.StatusOK, data)
}
