package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"upvotesystem/database"
	"upvotesystem/helpers"
	"upvotesystem/models"
	"upvotesystem/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateUUID() models.UUID {
	return models.UUID(uuid.New().String())
}

func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(c)
		defer cancel()
		var requestData models.PostRequest
		if err := c.ShouldBindBodyWithJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"details": helpers.CreateValidationErrorMessages(err),
			})
			return
		}
		post := models.Post{
			ID:          GenerateUUID(),
			Title:       requestData.Title,
			Content:     requestData.Content,
			ByUser:      requestData.ByUser,
			CreatedDate: helpers.RFC3339CurrentTime(),
			VotesCount:  0,
		}
		mappedPost, err := post.ToMap()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": err.Error(),
			})
			return
		}
		_, err = database.Client.HSet(ctx, "posts."+string(post.ID), mappedPost).Result()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"details": "post created",
			"post":    post,
		})

	}
}

func GetPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(c)
		defer cancel()
		var post models.Post
		post_id := c.Param("post_id")
		err := database.Client.HGetAll(ctx, "posts."+post_id).Scan(&post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": err.Error(),
			})
			return
		}
		if post.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"details": "post not found",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"post": post,
		})

	}
}

func GetPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var posts []models.Post
		var cursor uint64
		var keys []string
		var err error

		var minVotes int
		var filteredPosts []models.Post

		min_vote := c.Query("min_vote")
		if min_vote != "" {
			minVotes, err = strconv.Atoi(min_vote)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"details": "min_vote should be integer",
				})
				return
			}
			filteredPosts, err = services.CacheServiceInstance.GetFromCaches(c.Request.Context(), minVotes)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"details": "can not get from cache",
				})
				return
			}
			if filteredPosts != nil {
				c.JSON(http.StatusOK, gin.H{
					"posts":      filteredPosts,
					"count":      len(filteredPosts),
					"from_cache": true,
				})
				return
			}

		}
		for {
			keys, cursor, err = database.Client.Scan(c.Request.Context(), cursor, "posts.*", 100).Result()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			for _, key := range keys {
				var post models.Post
				database.Client.HGetAll(c.Request.Context(), key).Scan(&post)
				if post.Content != "" {
					posts = append(posts, post)
				}
			}

			if cursor == 0 {
				break
			}
		}
		if min_vote == "" {
			filteredPosts = posts

		} else {
			for _, item := range posts {
				if item.VotesCount > minVotes {
					filteredPosts = append(filteredPosts, item)
				}
			}

			err = services.CacheServiceInstance.SetCache(c.Request.Context(), minVotes, filteredPosts)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"details": "can not set cache",
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"posts": filteredPosts,
			"count": len(filteredPosts),
		})
	}
}

func UpVotePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(c)
		defer cancel()
		var requestData models.UpvoteRequest
		var post models.Post
		if err := c.ShouldBindBodyWithJSON(&requestData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": helpers.CreateValidationErrorMessages(err),
			})
			return
		}
		err := database.Client.HGetAll(ctx, "posts."+string(requestData.PostID)).Scan(&post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": err.Error(),
			})
			return
		}
		if post.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"details": "post not found",
			})
			return
		}
		postService := services.PostService{
			Post: &post,
		}
		newPost, vote, err := postService.UpVote(ctx, requestData.ByUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": err.Error(),
			})
			return
		}
		if newPost == nil {
			c.JSON(http.StatusOK, gin.H{
				"details": "Already upvoted",
			})
			return
		}
		services.CacheServiceInstance.InvalidateCache(c.Request.Context(), newPost.VotesCount)
		c.JSON(http.StatusOK, gin.H{
			"new_post": newPost,
			"vote":     vote,
		})

	}
}

func DownVotePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(c)
		defer cancel()
		var requestData models.UpvoteRequest
		var post models.Post
		if err := c.ShouldBindBodyWithJSON(&requestData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": helpers.CreateValidationErrorMessages(err),
			})
			return
		}
		err := database.Client.HGetAll(ctx, "posts."+string(requestData.PostID)).Scan(&post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": err.Error(),
			})
			return
		}
		if post.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"details": "post not found",
			})
			return
		}
		postService := services.PostService{
			Post: &post,
		}
		newPost, vote, err := postService.DownVote(ctx, requestData.ByUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": err.Error(),
			})
			return
		}
		if newPost == nil {
			c.JSON(http.StatusOK, gin.H{
				"details": "Already downvoted",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"new_post": newPost,
			"vote":     vote,
		})

	}
}
