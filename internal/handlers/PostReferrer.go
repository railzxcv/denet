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

func PostReferrer(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// id check
		id := chi.URLParam(r, "id")
		_, err := uuid.Parse(id)
		if err != nil {
			respondWithError(log, w, r, http.StatusBadRequest, errx.InvalidUserID, err)
			return
		}
		// body check
		var req requests.PostReferrer
		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			respondWithError(log, w, r, http.StatusBadRequest, errx.EmptyRequestBody, err)
			return
		}
		if req.ReferralCode == "" {
			respondWithError(log, w, r, http.StatusBadRequest, errx.InvalidReferralCode, nil)
			return
		}
		// referralCode - just an id of reffered user (referredBy)
		// TODO: check real referral code
		referredBy, err := uuid.Parse(req.ReferralCode)
		if err != nil {
			respondWithError(log, w, r, http.StatusBadRequest, errx.IncorrectReferralCode, err)
			return
		}

		err = storage.PostReferrer(id, req.ReferralCode, referredBy.String())
		if errors.Is(err, errx.ErrReferrerNotFound) {
			respondWithError(log, w, r, http.StatusBadRequest, errx.ReferrerNotFound, err)
			return
		}
		if errors.Is(err, errx.ErrUserNotFound) {
			respondWithError(log, w, r, http.StatusBadRequest, errx.UserNotFound, err)
			return
		}
		if err != nil {
			respondWithError(log, w, r, http.StatusInternalServerError, "failed to add referalCode", err)
			return
		}
		render.Status(r, http.StatusCreated)
	}
}
