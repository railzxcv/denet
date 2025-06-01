package handlers

import (
	"denet-app/internal/domain/users"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func GetUserStatus(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		_, err := uuid.Parse(id)
		if err != nil {
			respondWithError(log, w, r, http.StatusBadRequest, "invalid user ID format", err)
			return
		}
		var user users.StatusInfo

		user, err = storage.GetUserStatus(id)
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(log, w, r, http.StatusNotFound, "user not found", err)
			return
		}
		if err != nil {
			respondWithError(log, w, r, http.StatusInternalServerError, "failed to get user", err)
			return
		}
		fmt.Println(user)
		render.JSON(w, r, user)
	}
}
