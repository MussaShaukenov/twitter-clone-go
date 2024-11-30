package user

import (
	ctrl "MussaShaukenov/twitter-clone-go/internal/user/controller"
	"MussaShaukenov/twitter-clone-go/internal/user/repository"
	"MussaShaukenov/twitter-clone-go/internal/user/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"net/http"
)

func InitializeUserApp(db *pgxpool.Pool, redisClient *redis.Client, router chi.Router) (http.Handler, error) {
	dbRepo := repository.NewPostgres(db)
	redisRepo := repository.NewRedis(redisClient)

	uc := usecase.NewUserUseCase(dbRepo, redisRepo)
	controller := ctrl.NewController(uc)

	router.Mount("/users", ctrl.RegisterUserRoutes(controller))

	return router, nil
}
