package tweet

import (
	tweetCtrl "MussaShaukenov/twitter-clone-go/internal/tweet/controller"
	tweetRepo "MussaShaukenov/twitter-clone-go/internal/tweet/repository"
	tweetUc "MussaShaukenov/twitter-clone-go/internal/tweet/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func InitializeTweetApp(db *pgxpool.Pool, router chi.Router) (http.Handler, error) {
	repo := tweetRepo.NewPostgres(db)
	uc := tweetUc.NewTweetUseCase(repo)
	controller := tweetCtrl.NewController(uc)

	router.Mount("/tweets", tweetCtrl.RegisterTweetRoutes(controller))

	return router, nil
}
