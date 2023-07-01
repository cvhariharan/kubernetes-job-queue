package jobqueue

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

const JOB_QUEUE_NAME = "jobqueue"

type JobQueue interface {
	Publish(context.Context, chan string)
	Subscribe(context.Context, bool) chan string
}

type RedisQueue struct {
	client *redis.Client
}

func NewRedisJobqueue() JobQueue {
	opt, err := redis.ParseURL(os.Getenv("REDIS_CONNECTION_URL"))
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.NewClient(opt)
	return &RedisQueue{
		client: rdb,
	}
}

func (r *RedisQueue) Publish(ctx context.Context, jobsChan chan string) {
	for {
		job := <-jobsChan
		err := r.client.LPush(ctx, JOB_QUEUE_NAME, job).Err()
		if err != nil {
			log.Println(err)
		}
	}
}

func (r *RedisQueue) Subscribe(ctx context.Context, closeOnEmpty bool) chan string {
	jobsChan := make(chan string)
	go func(client *redis.Client, jobsChan chan string) {
		defer close(jobsChan)
		for {
			queueLength, err := client.LLen(ctx, JOB_QUEUE_NAME).Result()
			if err != nil {
				log.Println(err)
			}

			if closeOnEmpty && queueLength == 0 {
				return
			}

			job, err := client.BRPop(ctx, 0, JOB_QUEUE_NAME).Result()
			if err != nil {
				log.Println(err)
			}

			if len(job) > 1 {
				jobsChan <- job[1]
			}
		}
	}(r.client, jobsChan)

	return jobsChan
}
