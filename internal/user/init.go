package user

import (
	ctrl "MussaShaukenov/twitter-clone-go/internal/user/controller"
	"MussaShaukenov/twitter-clone-go/internal/user/repository"
	"MussaShaukenov/twitter-clone-go/internal/user/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func InitializeUserApp(db *pgxpool.Pool, router chi.Router) (http.Handler, error) {
	repo := repository.NewPostgres(db)
	uc := usecase.NewUserUseCase(repo)
	controller := ctrl.NewController(uc)

	router.Mount("/users", ctrl.RegisterUserRoutes(controller))

	return router, nil
}
