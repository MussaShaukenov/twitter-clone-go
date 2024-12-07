package followers

import "errors"

var (
	ErrInvalidID                       = errors.New("invalid id")
	ErrUserNotFound                    = errors.New("user not found")
	ErrIdsCannotBeTheSame              = errors.New("ids cannot be the same")
	ErrInvalidIDs                      = errors.New("invalid ids")
	ErrNotFollowing                    = errors.New("not following")
	ErrFollowerNotFound                = errors.New("follower not found")
	ErrFolloweeNotFound                = errors.New("followee not found")
	ErrFailedToCheckIfAlreadyFollowing = errors.New("failed to check if already following")
	ErrorAlreadyFollowing              = errors.New("already following")
)
