package avchecker

import (
	"bytes"
	"fmt"
	redis "github.com/xuyu/goredis"
	"net/http"
	"strings"
)

type Reporter interface {
	SendStats([]byte) error
	String() string
}

type redisQueueReporter struct {
	queueName string
	client    *redis.Redis
}

func (r *redisQueueReporter) SendStats(stats []byte) error {
	_, err := r.client.LPush(r.queueName, string(stats))
	return err
}

func (r *redisQueueReporter) String() string {
	return fmt.Sprintf("Redis 'LPUSH' to queue '%s'", r.queueName)
}

func urlWithDefaults(dialURL string) string {
	if !strings.Contains(dialURL, "?") {
		dialURL += "?"
	}
	if !strings.Contains(dialURL, "timeout=") {
		if !strings.HasSuffix(dialURL, "?") {
			dialURL += "&"
		}
		dialURL += "timeout=10s"
	}
	if !strings.Contains(dialURL, "maxidle=") {
		if !strings.HasSuffix(dialURL, "?") && !strings.HasSuffix(dialURL, "&") {
			dialURL += "&"
		}
		dialURL += "maxidle=1"
	}
	return dialURL
}

func NewRedisQueueReporter(queueName string, dialURL string) (Reporter, error) {
	client, err := redis.DialURL(urlWithDefaults(dialURL))
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

func (r *httpReporter) String() string {
	return fmt.Sprintf("http POST to %s", r.url)
}

func NewHttpReporter(url string, bodyType string) (Reporter, error) {
	return &httpReporter{
		client:   http.DefaultClient,
		bodyType: bodyType,
		url:      url,
	}, nil
}
