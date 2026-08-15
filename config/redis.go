package config

import (
	"context" //redis operation use a context
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client //this one creates a varibae that will hold our Redis client

func ConnectRedis() {
	addr := fmt.Sprintf("%s:%s",
	os.Getenv("REDIS_HOST"),
	os.Getenv("REDIS_PORT"),
    )

	RDB = redis.NewClient(&redis.Options{ //creating a redis client
		Addr: addr,
		DB: 0, //redis has 16 DBs 0-15, we use deafult 0 and also redis traditionally provides logical databases numbered:simply i haev chhosen a the deafault 0 
	})

	//ping redis to confirm connection
	if err := RDB.Ping(context.Background()).Err(); err!= nil {
		log.Fatal("failed to connect to redis:", err)
	}

	log.Println("Redis connected")
}

//.env  HOST and PORT
// go application creates redis:6379
//redis client :ping()
//redis container PONG