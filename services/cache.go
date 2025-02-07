package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"upvotesystem/database"
	"upvotesystem/models"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	CacheDuration time.Duration
}

func (cs *CacheService) SetCache(ctx context.Context, minVote int, result []models.Post) error {
	fmt.Printf("Setting cache minVote:%d\n", minVote)
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("top_posts_%d", minVote)
	_, err = database.Client.SetEx(ctx, cacheKey, data, cs.CacheDuration).Result()
	if err != nil {
		return err
	}
	err = database.Client.ZAdd(ctx, "cached_min_likes", redis.Z{
		Score:  float64(minVote),
		Member: minVote,
	}).Err()
	if err != nil {
		return err
	}
	return nil

}
func (cs *CacheService) InvalidateCache(ctx context.Context, newVotes int) error {
	keysToDelete, err := database.Client.ZRangeByScore(ctx, "cached_min_likes", &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%d", newVotes-1),
	}).Result()
	if err != nil {
		return err
	}
	if len(keysToDelete) > 0 {
		var cacheKeys []string
		for _, minVotes := range keysToDelete {
			fmt.Printf("cache for %s minVote should invalidate because we have a post with minVote %s now \n", minVotes, minVotes)
			cacheKeys = append(cacheKeys, fmt.Sprintf("top_posts_%s", minVotes))
		}
		err = database.Client.Del(ctx, cacheKeys...).Err()
		if err != nil {
			return err
		}
		err = database.Client.ZRem(ctx, "cached_min_likes", keysToDelete).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (cs *CacheService) GetFromCaches(ctx context.Context, minVotes int) ([]models.Post, error) {

	cacheKey := fmt.Sprintf("top_posts_%d", minVotes)

	cachedData, err := database.Client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var posts []models.Post
	err = json.Unmarshal([]byte(cachedData), &posts)
	fmt.Printf("get from cache minvote %d values = %v\n", minVotes, posts)
	if err != nil {
		return nil, err
	}
	return posts, nil

}

var CacheSerice CacheService = CacheService{
	CacheDuration: time.Second * 280,
}
