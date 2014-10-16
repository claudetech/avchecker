package avchecker

import (
	"net/http"
	"time"
)

type Checker struct {
	stats      *stats
	lastReport time.Time
	options    *Options
	url        string
	reporter   Reporter
	running    bool
}

func (c *Checker) createRequest() (req *http.Request, err error) {
	if req, err = http.NewRequest(c.options.RequestType, c.url, c.options.RequestBody); err != nil {
		return
	}
	req.Header.Add("User-Agent", "AvailabilityBot (https://github.com/claudetech/avchecker)")
	return
}

func (c *Checker) sendStats() {
	c.stats.compute()
	serializedStats, err := c.options.Formatter.Format(c.stats)
	c.lastReport = time.Now()
	c.stats.reset()

	if err != nil {
		c.options.Logger.Errorf("Could not format stats: %s", err.Error())
		return
	}
	c.options.Logger.Info("Sending stats.")

	if err := c.reporter.SendStats(serializedStats); err != nil {
		c.options.Logger.Errorf("Could not send stats: %s", err.Error())
	}
}

func (c *Checker) sendRequest() {
	req, err := c.createRequest()
	if err == nil {
		res, err := c.options.HttpClient.Do(req)
		if err != nil {
			c.options.Logger.Error(err.Error())
		} else if res.StatusCode >= 200 && res.StatusCode < 300 {
			c.stats.SuccessCount += 1
		}
	}
}

func (c *Checker) StartChecking() {
	runTime := 0
	c.stats.reset()
	c.lastReport = time.Now()
	c.running = true
	for c.running && (c.options.totalRuns == -1 || runTime < c.options.totalRuns) {
		runTime += 1
		c.stats.TryCount += 1
		c.sendRequest()
		if time.Now().Sub(c.lastReport) >= c.options.ReportInterval {
			c.sendStats()
		}
		time.Sleep(c.options.CheckInterval)
	}
	if c.stats.TryCount > 0 {
		c.sendStats()
	}
}

func (c *Checker) Stop() {
	c.running = false
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
