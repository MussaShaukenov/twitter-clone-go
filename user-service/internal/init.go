package user

import (
	ctrl "MussaShaukenov/twitter-clone-go/user-service/internal/controller"
	followerCtrl "MussaShaukenov/twitter-clone-go/user-service/internal/controller/followers"
	userCtrl "MussaShaukenov/twitter-clone-go/user-service/internal/controller/users"
	followerRepo "MussaShaukenov/twitter-clone-go/user-service/internal/repository/followers"
	otpRepo "MussaShaukenov/twitter-clone-go/user-service/internal/repository/otp"
	userRepo "MussaShaukenov/twitter-clone-go/user-service/internal/repository/users"
	followerUC "MussaShaukenov/twitter-clone-go/user-service/internal/usecase/followers"
	userUC "MussaShaukenov/twitter-clone-go/user-service/internal/usecase/users"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"net/http"
)

type Config struct {
	Db     *pgxpool.Pool
	Logger *zap.SugaredLogger
	Redis  *redis.Client
	Router *chi.Mux
}

func InitializeUserApp(config *Config) (http.Handler, error) {
	// initialize repositories
	userRepository := userRepo.NewUsersRepo(config.Db, config.Logger)
	followerRepository := followerRepo.NewFollowersRepo(config.Db, config.Logger)
	otpRepository := otpRepo.NewOTPRepo(config.Redis, config.Logger)

	// initialize use cases
	userUseCase := userUC.NewUserUseCase(userRepository, otpRepository, config.Logger)
	followerUseCase := followerUC.NewFollowerUseCase(userRepository, followerRepository, config.Logger)

	// initialize controller
	followerController := followerCtrl.NewFollowerController(followerUseCase, config.Logger)
	userController := userCtrl.NewUserController(userUseCase, config.Logger)

	// register routes
	config.Router.Mount("/users", ctrl.RegisterUserRoutes(userController))
	config.Router.Mount("/followers", ctrl.RegisterFollowerRoutes(followerController))

	return config.Router, nil
}
