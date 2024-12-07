package tweet

import (
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/controller"
	statsCtrl "MussaShaukenov/twitter-clone-go/tweet-service/internal/controller/stats"
	tagCtrl "MussaShaukenov/twitter-clone-go/tweet-service/internal/controller/tags"
	tweetCtrl "MussaShaukenov/twitter-clone-go/tweet-service/internal/controller/tweets"
	statsRepo "MussaShaukenov/twitter-clone-go/tweet-service/internal/repository/stats"
	tagRepo "MussaShaukenov/twitter-clone-go/tweet-service/internal/repository/tags"
	tweetRepo "MussaShaukenov/twitter-clone-go/tweet-service/internal/repository/tweets"
	statsUc "MussaShaukenov/twitter-clone-go/tweet-service/internal/usecase/stats"
	tagUc "MussaShaukenov/twitter-clone-go/tweet-service/internal/usecase/tags"
	tweetUc "MussaShaukenov/twitter-clone-go/tweet-service/internal/usecase/tweets"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Postgres *pgxpool.Pool
	Redis    *redis.Client
	Logger   *zap.SugaredLogger
	Router   *chi.Mux
	Mongo    *mongo.Database
}

func InitializeTweetApp(config *Config) (http.Handler, error) {
	tweetRepository := tweetRepo.NewTweetRepository(config.Postgres, config.Redis, 10*time.Minute)
	tagsRepository := tagRepo.NewTagsRepository(config.Postgres)
	statsRepository := statsRepo.NewTweetStatsRepository(config.Mongo)

	tweetUseCase := tweetUc.NewTweetUseCase(tweetRepository)
	tagsUseCase := tagUc.NewTagsUseCase(tagsRepository)
	statsUseCase := statsUc.NewTweetStatsUseCase(statsRepository)

	tweetController := tweetCtrl.NewController(tweetUseCase)
	tagsController := tagCtrl.NewTweetTagsController(tagsUseCase)
	statsController := statsCtrl.NewTweetStatsController(statsUseCase)

	config.Router.Mount("/tweets", controller.RegisterTweetRoutes(tweetController))
	config.Router.Mount("/tweets/tags", controller.RegisterTagsRoutes(tagsController))
	config.Router.Mount("/tweets/stats", controller.RegisterStatsRoutes(statsController))

	return config.Router, nil
}
