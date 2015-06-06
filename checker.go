package avchecker

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Checker struct {
	stats       *stats
	lastPublish time.Time
	options     *Options
	url         string
	reporter    Reporter
	running     bool
}

func (c *Checker) createRequest() (req *http.Request, err error) {
	if req, err = http.NewRequest(c.options.RequestType, c.url, c.options.RequestBody); err != nil {
		return
	}
	req.Header.Add("User-Agent", "AvailabilityBot (https://github.com/claudetech/avchecker)")
	for header, value := range c.options.RequestHeaders {
		req.Header.Add(header, value)
	}

	return
}

func (c *Checker) sendStats() {
	serializedStats, err := c.options.Formatter.Format(c.stats)
	c.lastPublish = time.Now()
	c.stats.reset()

	if err != nil {
		c.options.Logger.Warningf("Could not format stats: %s", err.Error())
		return
	}
	c.options.Logger.Infof("Sending stats via %s.", c.reporter.String())

	if err := c.reporter.SendStats(serializedStats); err != nil {
		c.options.Logger.Warningf("Could not send stats: %s", err.Error())
	}
}

func (c *Checker) sendRequest(drop bool) {
	req, err := c.createRequest()
	if err == nil {
		c.options.Logger.Tracef("start sending request to %s", req.URL.String())
		start := time.Now()
		res, err := c.options.HttpClient.Do(req)
		elapsed := time.Since(start)
		c.options.Logger.Debugf("request to %s sent in %dμs", req.URL.String(), elapsed.Nanoseconds()/1000)
		if drop {
			return
		}
		if err != nil {
			c.options.Logger.Warningf("Error during HTTP request: %s", err.Error())
		} else if res.StatusCode >= 200 && res.StatusCode < 300 {
			c.stats.SuccessCount += 1
			c.stats.totalTime += elapsed.Nanoseconds()
		}
		if res == nil || res.Body == nil {
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			c.options.Logger.Warningf("could not read received body: %s", err.Error())
		} else {
			c.options.Logger.Tracef("received body %s", body)
		}
	}
}

func (c *Checker) checkStats() {
	if c.stats.SuccessRatio < c.options.FatalThreshold {
		c.options.Logger.Fatalf("%s availability is under threshold. success ratio: %f", c.url, c.stats.SuccessRatio)
	}
}

func (c *Checker) StartChecking() {
	runTime := 0
	c.stats.reset()
	c.lastPublish = time.Now()
	c.running = true
	c.sendRequest(true) // warming up
	time.Sleep(c.options.CheckInterval)
	for c.running && (c.options.totalRuns == -1 || runTime < c.options.totalRuns) {
		runTime += 1
		c.stats.TryCount += 1
		c.sendRequest(false)
		if time.Now().Sub(c.lastPublish) >= c.options.PublishInterval {
			c.stats.compute()
			c.checkStats()
			c.sendStats()
		}
		c.options.Logger.Tracef("sleeping for %dms", c.options.CheckInterval.Nanoseconds()/1000000)
		time.Sleep(c.options.CheckInterval)
	}
	if c.stats.TryCount > 0 {
		c.sendStats()
	}
}

func (c *Checker) Stop() {
	c.running = false
}

func New(url string, reporter Reporter, options *Options) *Checker {
	return NewChecker(url, reporter, options)
}

func NewChecker(url string, reporter Reporter, options *Options) *Checker {
	options.setDefaults()
	return &Checker{
		stats:    &stats{extraFields: options.ExtraFields},
		url:      url,
		reporter: reporter,
		options:  options,
	}
}
