package stats

import (
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/domain"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	collection *mongo.Collection
}

func NewTweetStatsRepository(db *mongo.Database) *repository {
	return &repository{
		collection: db.Collection("stats"),
	}
}

func (repo *repository) GetTweetStats(ctx context.Context, tweetID int64) (*domain.TweetStats, error) {
	var stats domain.TweetStats
	err := repo.collection.FindOne(ctx, bson.M{"tweet_id": tweetID}).Decode(&stats)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If no stats exist, initialize with zero stats
			stats = domain.TweetStats{
				TweetID:  tweetID,
				Likes:    0,
				Dislikes: 0,
			}
			_, err = repo.collection.InsertOne(ctx, stats)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &stats, nil
}

func (repo *repository) UpdateLikes(ctx context.Context, tweetID int64, likesChange int64) error {
	_, err := repo.collection.UpdateOne(
		ctx,
		bson.M{"tweet_id": tweetID},
		bson.M{"$inc": bson.M{"likes": likesChange}},
	)
	return err
}

func (repo *repository) UpdateDislikes(ctx context.Context, tweetID int64, dislikesChange int64) error {
	_, err := repo.collection.UpdateOne(
		ctx,
		bson.M{"tweet_id": tweetID},
		bson.M{"$inc": bson.M{"dislikes": dislikesChange}},
	)
	return err
}
