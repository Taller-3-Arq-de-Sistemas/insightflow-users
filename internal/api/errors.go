package api

import (
	"encoding/json"
	"net/http"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func ErrorResponse(w http.ResponseWriter, status int, problemType, title, detail string) {
	prob := Problem{
		Type:   problemType,
		Title:  cases.Title(language.English).String(title),
		Status: status,
		Detail: cases.Title(language.English).String(detail),
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(prob)
}

func InternalServerError(w http.ResponseWriter, err error) {
	ErrorResponse(w, http.StatusInternalServerError, "https://statuscodes.io/500", "Internal Server Error", err.Error())
}

func BadRequest(w http.ResponseWriter, detail string) {
	ErrorResponse(w, http.StatusBadRequest, "https://statuscodes.io/400", "Bad Request", detail)
}

func NotFound(w http.ResponseWriter, detail string) {
	ErrorResponse(w, http.StatusNotFound, "https://statuscodes.io/404", "Not Found", detail)
}

func Unauthorized(w http.ResponseWriter, detail string) {
	ErrorResponse(w, http.StatusUnauthorized, "https://statuscodes.io/401", "Unauthorized", detail)
}

func Forbidden(w http.ResponseWriter, detail string) {
	ErrorResponse(w, http.StatusForbidden, "https://statuscodes.io/403", "Forbidden", detail)
}

func InvalidCredentials(w http.ResponseWriter, detail string) {
	ErrorResponse(w, http.StatusUnauthorized, "https://statuscodes.io/401", "Invalid Credentials", detail)
}

func Conflict(w http.ResponseWriter, detail string) {
	ErrorResponse(w, http.StatusConflict, "https://statuscodes.io/409", "Conflict", detail)
}

func ValidationErrors(w http.ResponseWriter, errors map[string]string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusUnprocessableEntity)

	prob := struct {
		Type   string            `json:"type"`
		Title  string            `json:"title"`
		Status int               `json:"status"`
		Errors map[string]string `json:"errors"`
	}{
		Type:   "https://statuscodes.io/422",
		Title:  "Validation failed",
		Status: http.StatusUnprocessableEntity,
		Errors: errors,
	}

	json.NewEncoder(w).Encode(prob)
}
