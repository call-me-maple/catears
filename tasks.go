package main

import (
	"context"
	"github.com/go-redis/redis/v7"
	"log"
)

func cancelTask(key string) error {
	value, err := client.Get(key).Result()
	switch err {
	case nil:
		_, err = messageSend.Cancel(context.Background(), value)
		if err != nil {
			log.Println(err)
		}
	case redis.Nil:
		return nil
	default:
		log.Println("error getting key: ", key, err)
		return err
	}
	return nil
}
