package usecase

import (
	"MussaShaukenov/twitter-clone-go/internal/tweet/domain"
	"MussaShaukenov/twitter-clone-go/internal/tweet/dto"
	"MussaShaukenov/twitter-clone-go/internal/tweet/repository"
	"errors"
	"fmt"
	"log"
)

type TweetUseCase interface {
	Create(dto dto.TweetDto) error
	Get(id int64) (*dto.TweetDto, error)
	List() ([]*dto.TweetDto, error)
	Update(in dto.TweetDto) (*dto.GetTweetResponse, error)
	Delete(id int) error
	GetUserTweets(id int) ([]*dto.TweetDto, error)
}

type tweetUseCase struct {
	tweetRepository repository.TweetRepository
}

func NewTweetUseCase(tweetRepository repository.TweetRepository) *tweetUseCase {
	return &tweetUseCase{
		tweetRepository: tweetRepository,
	}
}

func (uc *tweetUseCase) Create(dto dto.TweetDto) error {
	// Validation
	if len(dto.Title) == 0 {
		return errors.New("title cannot be empty")
	}
	if len(dto.Content) == 0 {
		return errors.New("content cannot be empty")
	}

	tweet := domain.ConvertFromDto(0, dto.Title, dto.Content, dto.Topic, dto.UserId)
	err := uc.tweetRepository.Insert(tweet)
	if err != nil {
		return err
	}
	return nil
}

// Validation
func (uc *tweetUseCase) Get(id int64) (*dto.TweetDto, error) {
	if id < 1 {
		return nil, fmt.Errorf("invalid ID: %v", id)
	}

	tweet, err := uc.tweetRepository.Get(id)
	if err != nil {
		log.Println("could not get a tweet")
		return nil, err
	}
	return domain.ConvertToDto(tweet), nil
}

func (uc *tweetUseCase) List() ([]*dto.TweetDto, error) {
	tweets, err := uc.tweetRepository.List()
	if err != nil {
		log.Println("could not list tweets")
		return nil, err
	}
	var result []*dto.TweetDto
	for _, tweet := range tweets {
		result = append(result, domain.ConvertToDto(tweet))
	}
	return result, nil
}

func (uc *tweetUseCase) Update(in dto.TweetDto) (*dto.GetTweetResponse, error) {
	tweet := domain.ConvertFromDto(in.ID, in.Title, in.Content, in.Topic, in.UserId)
	updatedTweet, err := uc.tweetRepository.Update(tweet)
	if err != nil {
		log.Println("could not update the updatedTweet")
		return nil, err
	}

	return domain.ConvertToGetTweetResponseDto(updatedTweet), nil
}

func (uc *tweetUseCase) Delete(id int) error {
	if id < 1 {
		return fmt.Errorf("invalid ID: %v", id)
	}

	err := uc.tweetRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("could not delete: %w", err)
	}
	return nil
}

func (uc *tweetUseCase) GetUserTweets(id int) ([]*dto.TweetDto, error) {
	if id < 1 {
		fmt.Errorf("invalid ID: %v", id)
	}
	tweets, err := uc.tweetRepository.GetUserTweets(id)
	if err != nil {
		log.Printf("could not get tweets of user %v", id)
		return nil, err
	}
	var result []*dto.TweetDto
	for _, tweet := range tweets {
		result = append(result, domain.ConvertToDto(tweet))
	}
	return result, nil
}
