package followers

import (
	"MussaShaukenov/twitter-clone-go/user-service/internal/domain"
	"MussaShaukenov/twitter-clone-go/user-service/internal/repository"
	"errors"
	"go.uber.org/zap"
)

type useCase struct {
	userRepo     repository.UserRepo
	followerRepo repository.FollowerRepo
	logger       *zap.SugaredLogger
}

func NewFollowerUseCase(userRepo repository.UserRepo, followerRepo repository.FollowerRepo, logger *zap.SugaredLogger) *useCase {
	return &useCase{
		userRepo:     userRepo,
		followerRepo: followerRepo,
		logger:       logger,
	}
}

func (uc *useCase) Follow(followerID, followeeID int) error {
	// validate ids
	err := validateIds(followerID, followeeID)
	if err != nil {
		uc.logger.Warnw("validation failed: validateIds", "followerID", followerID, "followeeID:", followeeID, "error", err)
		return err
	}
	// check if ids exists
	err = uc.checkIfFollowerAndFolloweeExist(followerID, followeeID)
	if err != nil {
		uc.logger.Warnw("validation failed: checkIfFollowerAndFolloweeExist",
			"followerID", followerID, "followeeID:", followeeID, "error", err)
		return err
	}

	// check if already following
	isFollowing, err := uc.followerRepo.IsFollowing(followerID, followeeID)
	if err != nil {
		uc.logger.Errorw("failed to check following status", "followerID", followerID, "followeeID", followeeID, "error", err)
		return ErrFailedToCheckIfAlreadyFollowing
	}
	if isFollowing {
		return ErrorAlreadyFollowing
	}

	// follow
	err = uc.followerRepo.Follow(followerID, followeeID)
	if err != nil {
		uc.logger.Errorw("Failed to follow", "followerID", followerID, "followeeID", followeeID, "error", err)
		return err
	}

	return nil
}

func (uc *useCase) Unfollow(followerID, followeeID int) error {
	// validate ids
	err := validateIds(followerID, followeeID)
	if err != nil {
		uc.logger.Warnw("validation failed: validateIds", "followerID", followerID, "followeeID:", followeeID, "error", err)
		return err
	}
	// check if ids exists
	err = uc.checkIfFollowerAndFolloweeExist(followerID, followeeID)
	if err != nil {
		uc.logger.Warnw("validation failed: checkIfFollowerAndFolloweeExist", "followerID", followerID, "followeeID:", followeeID, "error", err)
		return err
	}

	// check if already following
	isFollowing, err := uc.followerRepo.IsFollowing(followerID, followeeID)
	if err != nil {
		uc.logger.Errorw("failed to check following status", "followerID", followerID, "followeeID", followeeID, "error", err)
		return ErrFailedToCheckIfAlreadyFollowing
	}
	if !isFollowing {
		return ErrNotFollowing
	}

	// unfollow
	err = uc.followerRepo.Unfollow(followerID, followeeID)
	if err != nil {
		uc.logger.Errorw("Failed to unfollow", "followerID", followerID, "followeeID", followeeID, "error", err)
		return err
	}

	return nil
}

func (uc *useCase) GetFollowers(userID int) ([]*domain.User, error) {
	// validate id
	err := validateID(userID)
	if err != nil {
		uc.logger.Warnw("validation failed: validateID", "userID", userID, "error", err)
		return nil, err
	}
	// check if id exists
	_, err = uc.userRepo.GetByID(userID)
	if err != nil {
		uc.logger.Warnw("user not found in GetFollowers", "userID", userID, "error", err)
		switch {
		case errors.Is(err, ErrUserNotFound):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	followers, err := uc.followerRepo.GetFollowers(userID)
	if err != nil {
		uc.logger.Errorw("Failed to retrieve followers", "userID", userID, "error", err)
		return nil, err
	}

	return followers, nil
}

func (uc *useCase) GetFollowing(userID int) ([]*domain.User, error) {
	// validate id
	err := validateID(userID)
	if err != nil {
		uc.logger.Warnw("Validation failed in GetFollowing", "userID", userID, "error", err)
		return nil, err
	}

	// check if id exists
	_, err = uc.userRepo.GetByID(userID)
	if err != nil {
		uc.logger.Warnw("User not found in GetFollowing", "userID", userID, "error", err)
		switch {
		case errors.Is(err, ErrUserNotFound):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	following, err := uc.followerRepo.GetFollowing(userID)
	if err != nil {
		uc.logger.Errorw("Failed to retrieve following", "userID", userID, "error", err)
		return nil, err
	}

	return following, nil
}

func (uc *useCase) IsFollowing(followerID, followeeID int) (bool, error) {
	// validate ids
	err := validateIds(followerID, followeeID)
	if err != nil {
		uc.logger.Warnw("Validation failed in IsFollowing", "followerID", followerID, "followeeID", followeeID, "error", err)
		return false, err
	}
	// check if ids exists
	err = uc.checkIfFollowerAndFolloweeExist(followerID, followeeID)
	if err != nil {
		uc.logger.Warnw("Follower or followee does not exist in IsFollowing", "followerID", followerID, "followeeID", followeeID, "error", err)
		return false, err
	}

	isFollowing, err := uc.followerRepo.IsFollowing(followerID, followeeID)
	if err != nil {
		uc.logger.Errorw("Failed to check if following status in IsFollowing", "followerID", followerID, "followeeID", followeeID, "error", err)
		return false, err
	}
	return isFollowing, nil
}

func validateIds(followerID, followeeID int) error {
	if followerID == followeeID {
		return ErrIdsCannotBeTheSame
	}
	if followerID < 1 || followeeID < 1 {
		return ErrInvalidIDs
	}
	return nil
}

func (uc *useCase) checkIfFollowerExists(followerID int) error {
	if _, err := uc.userRepo.GetByID(followerID); err != nil {
		return ErrFollowerNotFound
	}
	return nil
}

func (uc *useCase) checkIfFolloweeExists(followeeID int) error {
	if _, err := uc.userRepo.GetByID(followeeID); err != nil {
		return ErrFolloweeNotFound
	}
	return nil
}

func (uc *useCase) checkIfFollowerAndFolloweeExist(followerID, followeeID int) error {
	followerErr := uc.checkIfFollowerExists(followerID)
	followeeErr := uc.checkIfFolloweeExists(followeeID)

	if followerErr != nil && followeeErr != nil {
		return errors.Join(followeeErr, followeeErr)
	}
	if followerErr != nil && followeeErr == nil {
		return followerErr
	}
	if followerErr == nil && followeeErr != nil {
		return followeeErr
	}
	return nil
}

func validateID(id int) error {
	if id < 1 {
		return ErrInvalidID
	}
	return nil
}
