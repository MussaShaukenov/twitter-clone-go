package tags

import (
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/dto"
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/repository"
	"fmt"
)

type tweetUseCase struct {
	tweetTagRepository repository.TweetTagRepository
}

func NewTagsUseCase(tweetTagRepository repository.TweetTagRepository) *tweetUseCase {
	return &tweetUseCase{
		tweetTagRepository: tweetTagRepository,
	}
}

func (uc *tweetUseCase) AddTag(tweetId int64, tagId int64) error {
	if tweetId < 1 {
		return fmt.Errorf("invalid tweetId: %v", tweetId)
	}
	if tagId < 1 {
		return fmt.Errorf("invalid tagId: %v", tagId)
	}

	err := uc.tweetTagRepository.AddTag(tweetId, tagId)
	if err != nil {
		return fmt.Errorf("could not add tag to tweet: %w", err)
	}
	return nil
}

func (uc *tweetUseCase) GetTweetTags(tweetId int64) ([]*dto.TagDto, error) {
	if tweetId < 1 {
		return nil, fmt.Errorf("invalid tweetId: %v", tweetId)
	}

	tags, err := uc.tweetTagRepository.GetTweetTags(tweetId)
	if err != nil {
		return nil, fmt.Errorf("could not get tags of tweet: %w", err)
	}

	var result []*dto.TagDto
	for _, tag := range tags {
		result = append(result, &dto.TagDto{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	return result, nil
}

func (uc *tweetUseCase) ListTags() ([]*dto.TagDto, error) {
	tags, err := uc.tweetTagRepository.ListTags()
	if err != nil {
		return nil, fmt.Errorf("could not list tags: %w", err)
	}

	var result []*dto.TagDto
	for _, tag := range tags {
		result = append(result, &dto.TagDto{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	return result, nil
}
