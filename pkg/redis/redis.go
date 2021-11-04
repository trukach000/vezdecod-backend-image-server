package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func InitRedisClient(
	host string,
	port string,
	password string,
	dbName string,
	establishmentRetryPeriod string,
) (*redis.Client, error) {
	var err error

	database, err := strconv.Atoi(dbName)
	if err != nil {
		return nil, err
	}

	connectOptions := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       database,
	}

	err = waitRedisReachable(connectOptions, establishmentRetryPeriod, 2*time.Minute)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(connectOptions), nil
}

func tryConnectToRedis(connectOptions *redis.Options) bool {
	redisClient := redis.NewClient(connectOptions)
	ctx := context.Background()
	defer ctx.Done()

	err := redisClient.Ping().Err()
	if err != nil {
		logrus.Warningf("can't ping redis: %s", err)
		return false
	}

	_ = redisClient.Close()

	return true
}

func waitRedisReachable(connectOptions *redis.Options, establishmentRetryPeriod string, maxWait time.Duration) error {
	failed := time.Now().Add(maxWait)
	retryPeriod, err := time.ParseDuration(establishmentRetryPeriod)
	if err != nil {
		return err
	}

	for time.Now().Before(failed) {
		if tryConnectToRedis(connectOptions) {
			logrus.Infof("redis connection established")
			return nil
		}
		time.Sleep(retryPeriod)
	}

	return errors.New(fmt.Sprintf("can't establish connection with given options: %#v", connectOptions))
}
