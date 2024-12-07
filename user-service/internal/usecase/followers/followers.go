package followers

import (
	"MussaShaukenov/twitter-clone-go/user-service/internal/domain"
	"MussaShaukenov/twitter-clone-go/user-service/internal/repository"
	"errors"
)

type useCase struct {
	userRepo     repository.UserRepo
	followerRepo repository.FollowerRepo
}

func NewFollowerUseCase(userRepo repository.UserRepo, followerRepo repository.FollowerRepo) *useCase {
	return &useCase{
		userRepo:     userRepo,
		followerRepo: followerRepo,
	}
}

func (uc *useCase) Follow(followerID, followeeID int) error {
	// validate ids
	err := validateIds(followerID, followeeID)
	if err != nil {
		return err
	}
	// check if ids exists
	err = uc.checkIfFollowerAndFolloweeExist(followerID, followeeID)
	if err != nil {
		return err
	}

	// check if already following
	isFollowing, err := uc.followerRepo.IsFollowing(followerID, followeeID)
	if err != nil {
		return ErrFailedToCheckIfAlreadyFollowing
	}
	if isFollowing {
		return ErrorAlreadyFollowing
	}

	// follow
	return uc.followerRepo.Follow(followerID, followeeID)
}

func (uc *useCase) Unfollow(followerID, followeeID int) error {
	// validate ids
	err := validateIds(followerID, followeeID)
	if err != nil {
		return err
	}
	// check if ids exists
	err = uc.checkIfFollowerAndFolloweeExist(followerID, followeeID)
	if err != nil {
		return err
	}

	// check if already following
	isFollowing, err := uc.followerRepo.IsFollowing(followerID, followeeID)
	if err != nil {
		return ErrFailedToCheckIfAlreadyFollowing
	}
	if !isFollowing {
		return ErrNotFollowing
	}

	// unfollow
	return uc.followerRepo.Unfollow(followerID, followeeID)
}

func (uc *useCase) GetFollowers(userID int) ([]*domain.User, error) {
	// validate id
	err := validateID(userID)
	if err != nil {
		return nil, err
	}
	// check if id exists
	_, err = uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return uc.followerRepo.GetFollowers(userID)
}

func (uc *useCase) GetFollowing(userID int) ([]*domain.User, error) {
	// validate id
	err := validateID(userID)
	if err != nil {
		return nil, err
	}

	// check if id exists
	_, err = uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return uc.followerRepo.GetFollowing(userID)
}

func (uc *useCase) IsFollowing(followerID, followeeID int) (bool, error) {
	// validate ids
	err := validateIds(followerID, followeeID)
	if err != nil {
		return false, err
	}
	// check if ids exists
	err = uc.checkIfFollowerAndFolloweeExist(followerID, followeeID)
	if err != nil {
		return false, err
	}

	return uc.followerRepo.IsFollowing(followerID, followeeID)
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
