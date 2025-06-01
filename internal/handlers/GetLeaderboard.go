package handlers

import (
	"denet-app/internal/domain/users"
	"log/slog"
	"net/http"
	"denet-app/internal/errx"
	"github.com/go-chi/render"
)

type Storage interface {
	GetLeaderboard() ([]users.LeaderboardEntry, error)
	GetUserStatus(id string) (users.StatusInfo, error)
	PostTaskComplete(id, taskType string) error
	PostReferrer(id, referralCode, referredBy string) error
}

func respondWithError(log *slog.Logger, w http.ResponseWriter, r *http.Request, statusCode int,  errorMsg string, err error) {
	log.Error(errorMsg, slog.Any("error", err))
	errorResponse := errx.ErrorResponse{Error: errorMsg}	
	render.Status(r, statusCode)
	render.JSON(w,r,errorResponse)
}


func GetLeaderboard(log *slog.Logger, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var leaders []users.LeaderboardEntry
		// TODO: add redis cache support
		leaders, err := storage.GetLeaderboard()
		if err != nil {
			respondWithError(log, w,r,  http.StatusInternalServerError, "failed to get leaderboard", err)
			return
		}
		render.JSON(w, r, leaders)
	}
}
