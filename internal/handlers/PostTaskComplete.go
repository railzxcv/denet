package handlers
import (
	"denet-app/internal/domain/requests"
	errx "denet-app/internal/errx"
	"log/slog"
	"net/http"

	"errors"
	"io"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func PostTaskComplete(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// id check
		id := chi.URLParam(r, "id")
		_, err := uuid.Parse(id)
		if err != nil {
			respondWithError(log, w, r, http.StatusBadRequest, errx.InvalidUserID, err)
			return
		}
		// body check
		var req requests.PostTaskComplete
		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			respondWithError(log, w, r, http.StatusBadRequest, errx.EmptyRequestBody, err)
			return
		}
		if req.TaskType == "" {
			respondWithError(log, w, r, http.StatusBadRequest, errx.InvalidTaskType, nil)
			return
		}

		err = storage.PostTaskComplete(id, req.TaskType)
		if errors.Is(err, errx.ErrUserNotFound) {
			respondWithError(log, w, r, http.StatusBadRequest, errx.UserNotFound, err)
			return
		}
		if errors.Is(err, errx.ErrTaskNotFound) {
			respondWithError(log, w, r, http.StatusBadRequest, errx.TaskNotFound, err)
			return
		}
		if errors.Is(err, errx.ErrNoChange) {
			render.NoContent(w, r)
			return
		}
		if err != nil {
			respondWithError(log, w, r, http.StatusInternalServerError, "failed to complete task", err)
			return
		}
		render.Status(r, http.StatusCreated)
	}
}