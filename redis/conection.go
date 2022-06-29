package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	config "github.com/spf13/viper"
)

var RedisClient = conection()

func conection() *redis.Client {

	config.SetConfigFile(".env")
	config.ReadInConfig()

	host := config.GetString("host")

	port := config.GetString("REDISport")

	addres := host + ":" + port

	client := redis.NewClient(&redis.Options{
		Addr:     addres,
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}

	fmt.Println(pong)

	return client
}

func Set(key string, value interface{}, expire time.Duration) error {
	p, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return RedisClient.Set(key, p, expire).Err()
}

func Get(key string) (interface{}, error) {

	p, err := RedisClient.Get(key).Result()

	return p, err
}
