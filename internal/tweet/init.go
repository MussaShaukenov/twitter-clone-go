package tweet

import (
	tweetCtrl "MussaShaukenov/twitter-clone-go/internal/tweet/controller"
	tweetRepo "MussaShaukenov/twitter-clone-go/internal/tweet/repository"
	tweetUc "MussaShaukenov/twitter-clone-go/internal/tweet/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

func InitializeTweetApp(db *pgxpool.Pool, redisClient *redis.Client, router chi.Router) (http.Handler, error) {
	repo := tweetRepo.NewPostgres(db, redisClient, 10*time.Minute)
	uc := tweetUc.NewTweetUseCase(repo)
	controller := tweetCtrl.NewController(uc)

	router.Mount("/tweets", tweetCtrl.RegisterTweetRoutes(controller))
	router.Mount("/tags", tweetCtrl.RegisterTagsRoutes(controller))

	return router, nil
}
