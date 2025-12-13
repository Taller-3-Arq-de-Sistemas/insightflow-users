package users

import (
	"net/http"

	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/api"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/auth"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/json"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/validator"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *svc
}

func NewHandler(service *svc) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var params CreateUserParams

	if err := json.Read(r, &params); err != nil {
		api.BadRequest(w, err.Error())
		return
	}

	if errors := validator.ValidateStruct(params); errors != nil {
		api.ValidationErrors(w, errors)
		return
	}

	user, err := h.service.CreateUser(r.Context(), params)
	if err != nil {
		if err == ErrInvalidStatus || err == ErrInvalidRole {
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

	json.Write(w, http.StatusCreated, user)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		api.InternalServerError(w, err)
		return
	}

	json.Write(w, http.StatusOK, users)
}

func (h *Handler) FindUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.service.FindUserById(r.Context(), id)
	if err != nil {
		if err == ErrUserNotFound {
			api.NotFound(w, err.Error())
			return
		}
		api.InternalServerError(w, err)
		return
	}

	json.Write(w, http.StatusOK, user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var params UpdateUserParams

	if err := json.Read(r, &params); err != nil {
		api.BadRequest(w, err.Error())
		return
	}

	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		api.Unauthorized(w, "unauthorized")
		return
	}

	if userID != id {
		api.Forbidden(w, "you can only update your own profile")
		return
	}

	user, err := h.service.UpdateUser(r.Context(), id, params)
	if err != nil {
		if err == ErrUserNotFound {
			api.NotFound(w, err.Error())
			return
		}
		if err == ErrUserAlreadyExists {
			api.Conflict(w, err.Error())
			return
		}
		if err == ErrInvalidFullname {
			api.BadRequest(w, err.Error())
			return
		}
		api.InternalServerError(w, err)
		return
	}

	json.Write(w, http.StatusOK, user)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.service.DeleteUser(r.Context(), id)
	if err != nil {
		if err == ErrUserNotFound {
			api.NotFound(w, err.Error())
			return
		}
		api.InternalServerError(w, err)
		return
	}

	json.Write(w, http.StatusOK, nil)
}
