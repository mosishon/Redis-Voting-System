package services

import (
	"context"
	"fmt"
	"strconv"
	"upvotesystem/database"
	"upvotesystem/models"
)

type PostService struct {
	Post *models.Post
}

func (service *PostService) UpVote(ctx context.Context, ByUser int) (*models.Post, *models.Vote, error) {
	var vote models.Vote
	err := database.Client.HGetAll(ctx, "votes."+string(service.Post.ID)+"."+strconv.Itoa(ByUser)).Scan(&vote)
	fmt.Println("getting: " + "votes." + string(service.Post.ID) + "." + strconv.Itoa(ByUser))
	if err != nil {
		return nil, nil, err
	}
	if vote.PostID == "" || vote.Type == models.DownVote {
		// Not voted or downvoted
		newVoteCount := service.Post.VotesCount + 1
		if vote.Type == models.DownVote {
			newVoteCount++
		}
		newPost := models.Post{
			ID:         service.Post.ID,
			Title:      service.Post.Title,
			Content:    service.Post.Content,
			ByUser:     service.Post.ByUser,
			VotesCount: newVoteCount,
		}
		mappedPost, err := newPost.ToMap()
		if err != nil {
			return nil, nil, err

		}
		_, err = database.Client.HSet(ctx, "posts."+string(service.Post.ID), mappedPost).Result()
		if err != nil {
			return nil, nil, err

		}
		vote = models.Vote{
			Type:   models.UpVote,
			ByUser: ByUser,
			PostID: service.Post.ID,
		}
		mappedVote, err := vote.ToMap()
		if err != nil {
			return nil, nil, err

		}
		_, err = database.Client.HSet(ctx, "votes."+string(service.Post.ID)+"."+strconv.Itoa(ByUser), mappedVote).Result()
		if err != nil {
			return nil, nil, err

		}
		return &newPost, &vote, nil
	} else {
		return nil, nil, nil
	}

}

func (service *PostService) DownVote(ctx context.Context, ByUser int) (*models.Post, *models.Vote, error) {
	var vote models.Vote
	err := database.Client.HGetAll(ctx, "votes."+string(service.Post.ID)+"."+strconv.Itoa(ByUser)).Scan(&vote)
	fmt.Println("getting: " + "votes." + string(service.Post.ID) + "." + strconv.Itoa(ByUser))
	if err != nil {
		return nil, nil, err
	}
	if vote.PostID == "" || vote.Type == models.UpVote {
		// Not voted or upvoted

		newVoteCount := service.Post.VotesCount - 1
		if vote.Type == models.DownVote {
			newVoteCount--
		}
		newPost := models.Post{
			ID:         service.Post.ID,
			Title:      service.Post.Title,
			Content:    service.Post.Content,
			ByUser:     service.Post.ByUser,
			VotesCount: newVoteCount,
		}
		mappedPost, err := newPost.ToMap()
		if err != nil {
			return nil, nil, err

		}
		_, err = database.Client.HSet(ctx, "posts."+string(service.Post.ID), mappedPost).Result()
		if err != nil {
			return nil, nil, err

		}
		vote = models.Vote{
			Type:   models.DownVote,
			ByUser: ByUser,
			PostID: service.Post.ID,
		}
		mappedVote, err := vote.ToMap()
		if err != nil {
			return nil, nil, err

		}
		_, err = database.Client.HSet(ctx, "votes."+string(service.Post.ID)+"."+strconv.Itoa(ByUser), mappedVote).Result()
		if err != nil {
			return nil, nil, err

		}
		return &newPost, &vote, nil
	} else {
		return nil, nil, nil
	}

}
