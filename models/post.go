package models

import (
	"encoding/json"
	"time"
)

type VoteType string

type UUID string

const (
	UpVote   VoteType = "UP"
	DownVote VoteType = "DOWN"
)

type Vote struct {
	Type   VoteType `json:"type" redis:"type"`
	ByUser int      `json:"by_user" redis:"by_user"`
	PostID UUID     `json:"post_id" redis:"post_id"`
}

// database model
type Post struct {
	ID          UUID      `json:"id" redis:"id"`
	Title       string    `json:"title" redis:"title"`
	Content     string    `json:"content" redis:"content"`
	ByUser      int       `json:"by_user" redis:"by_user"`
	VotesCount  int       `json:"votes_count" redis:"votes_count"`
	CreatedDate time.Time `json:"created_date" redis:"created_date"`
}

// Request data
type PostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	ByUser  int    `json:"by_user" binding:"required"`
}

type UpvoteRequest struct {
	ByUser int  `json:"by_user" redis:"by_user" binding:"required"`
	PostID UUID `json:"post_id" redis:"post_id" binding:"required"`
}

// use refrence to save memory
func (post *Post) ToMap() (map[string]interface{}, error) {
	data, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// use refrence to save memory

func (vote *Vote) ToMap() (map[string]interface{}, error) {
	data, err := json.Marshal(vote)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}
