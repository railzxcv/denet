package routes

import (
	mwAuth "denet-app/internal/middleware/auth"
	"denet-app/internal/storage"
	"denet-app/internal/token"
	"log/slog"
	"net/http"

	"denet-app/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(log *slog.Logger, storage storage.Storage, TokenMn *token.TokenManager) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	// выдача jwt-токенов идет с другого сервиса, поэтому пока что доступ без токена 
	router.Route("/users", func(r chi.Router) {
		r.Get("/{id}/status", handlers.GetUserStatus(log, storage))
		r.Get("/leaderboard", handlers.GetLeaderboard(log, storage))
		r.Post("/{id}/task/complete", handlers.PostTaskComplete(log, storage))
		r.Post("/{id}/referrer", handlers.PostReferrer(log, storage))
	})
	// for jwt-access
	router.Route("/vipusers", func(r chi.Router) {
		r.Use(mwAuth.AuthMiddleware(TokenMn, log))
		r.Get("/users/{id}/status", handlers.GetUserStatus(log, storage))
		r.Get("/leaderboard", handlers.GetLeaderboard(log, storage))
		r.Post("/{id}/task/complete", handlers.PostTaskComplete(log, storage))
		r.Post("/{id}/referrer", handlers.PostReferrer(log, storage))
	})

	return router
}
