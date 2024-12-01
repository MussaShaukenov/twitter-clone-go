package user

import (
	ctrl "MussaShaukenov/twitter-clone-go/user-service/internal/controller"
	repo "MussaShaukenov/twitter-clone-go/user-service/internal/repository"
	uc "MussaShaukenov/twitter-clone-go/user-service/internal/usecase"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitializeUserApp(db *pgxpool.Pool, redisClient *redis.Client, router chi.Router) (http.Handler, error) {
	dbRepo := repo.NewPostgres(db)
	redisRepo := repo.NewRedis(redisClient)

	uc := uc.NewUserUseCase(dbRepo, redisRepo)
	controller := ctrl.NewController(uc)

	router.Mount("/users", ctrl.RegisterUserRoutes(controller))

	return router, nil
}
