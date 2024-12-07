package stats

import (
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/domain"
	repo "MussaShaukenov/twitter-clone-go/tweet-service/internal/repository"
	"context"
)

type useCase struct {
	repo repo.TweetStatsRepo
}

func NewTweetStatsUseCase(repo repo.TweetStatsRepo) *useCase {
	return &useCase{
		repo: repo,
	}
}

func (uc *useCase) GetTweetStats(ctx context.Context, tweetID int64) (*domain.TweetStats, error) {
	tweetStats, err := uc.repo.GetTweetStats(ctx, tweetID)
	if err != nil {
		return nil, err
	}
	return tweetStats, nil
}

func (uc *useCase) AddLike(ctx context.Context, tweetID int64) error {
	err := uc.repo.UpdateLikes(ctx, tweetID, 1)
	if err != nil {
		return err
	}
	return nil
}

func (uc *useCase) AddDislike(ctx context.Context, tweetID int64) error {
	err := uc.repo.UpdateDislikes(ctx, tweetID, 1)
	if err != nil {
		return err
	}
	return nil
}

func (uc *useCase) RemoveLike(ctx context.Context, tweetID int64) error {
	err := uc.repo.UpdateLikes(ctx, tweetID, -1)
	if err != nil {
		return err
	}
	return nil
}

func (uc *useCase) RemoveDislike(ctx context.Context, tweetID int64) error {
	err := uc.repo.UpdateDislikes(ctx, tweetID, -1)
	if err != nil {
		return err
	}
	return nil
}
