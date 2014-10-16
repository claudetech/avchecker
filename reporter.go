package avchecker

import (
	redis "github.com/xuyu/goredis"
)

type Reporter interface {
	SendStats([]byte) error
}

type redisQueueReporter struct {
	queueName string
	client    *redis.Redis
}

func (r *redisQueueReporter) SendStats(stats []byte) error {
	_, err := r.client.LPush(r.queueName, string(stats))
	return err
}

func NewRedisQueueReporter(queueName string, options *redis.DialConfig) (Reporter, error) {
	client, err := redis.Dial(options)
	if err != nil {
		return nil, err
	}
	return &redisQueueReporter{
		queueName: queueName,
		client:    client,
	}, nil
}
