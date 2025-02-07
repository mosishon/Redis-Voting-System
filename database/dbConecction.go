package database

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func initdb() *redis.Client {
	opt, err := redis.ParseURL("redis://root:@localhost:6379/0")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Redis")

	client := redis.NewClient(opt)
	return client
}

var Client *redis.Client

func init() {
	Client = initdb()
	log.Println("Redis instance created")
}
