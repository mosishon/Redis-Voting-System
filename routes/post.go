package routes

import (
	"upvotesystem/controllers"

	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.RouterGroup) {
	router.POST("/posts", controllers.CreatePost())
	router.GET("/posts", controllers.GetPosts())
	router.GET("/post/:post_id", controllers.GetPost())
	router.POST("/post/upvote", controllers.UpVotePost())
	router.POST("/post/downvote", controllers.DownVotePost())

}
