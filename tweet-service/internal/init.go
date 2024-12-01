package tweet

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	ctrl "MussaShaukenov/twitter-clone-go/tweet-service/internal/controller"
	repo "MussaShaukenov/twitter-clone-go/tweet-service/internal/repository"
	uc "MussaShaukenov/twitter-clone-go/tweet-service/internal/usecase"
)

func InitializeTweetApp(db *pgxpool.Pool, redisClient *redis.Client, router chi.Router) (http.Handler, error) {
	repo := repo.NewPostgres(db, redisClient, 10*time.Minute)
	uc := uc.NewTweetUseCase(repo)
	controller := ctrl.NewController(uc)

	router.Mount("/tweets", ctrl.RegisterTweetRoutes(controller))
	router.Mount("/tags", ctrl.RegisterTagsRoutes(controller))

	return router, nil
}
