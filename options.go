package avchecker

import (
	"github.com/claudetech/loggo"
	"github.com/claudetech/loggo/default"
	"io"
	"net/http"
	"time"
)

type miniHttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Options struct {
	CheckInterval  time.Duration
	ReportInterval time.Duration
	Logger         *loggo.Logger
	Formatter      Formatter
	HttpClient     miniHttpClient
	RequestType    string
	RequestBody    io.Reader
	ExtraFields    map[string]interface{}
	totalRuns      int
}

func (o *Options) setDefaults() {
	if o.CheckInterval == 0 {
		o.CheckInterval = 10 * time.Second
	}
	if o.ReportInterval == 0 {
		o.ReportInterval = 1 * time.Minute
	}
	if o.Logger == nil {
		o.Logger = loggo_default.Log
	}
	if o.Formatter == nil {
		o.Formatter = &JsonFormatter{}
	}
	if o.ExtraFields == nil {
		o.ExtraFields = make(map[string]interface{})
	}
	if o.HttpClient == nil {
		o.HttpClient = http.DefaultClient
	}
	if o.RequestType == "" {
		o.RequestType = "GET"
	}
	if o.totalRuns == 0 {
		o.totalRuns = -1
	}
}
