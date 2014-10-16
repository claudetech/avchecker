package avchecker

import (
	"bytes"
	redis "github.com/xuyu/goredis"
	"net/http"
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

type httpReporter struct {
	url      string
	bodyType string
	client   *http.Client
}

func (h *httpReporter) SendStats(stats []byte) error {
	_, err := h.client.Post(h.url, h.bodyType, bytes.NewReader(stats))
	return err
}

func NewHttpReporter(url string, bodyType string) (Reporter, error) {
	return &httpReporter{
		client:   http.DefaultClient,
		bodyType: bodyType,
		url:      url,
	}, nil
}
