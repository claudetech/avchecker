package main

import (
	"bytes"
	"fmt"
	"github.com/claudetech/avchecker"
	"github.com/claudetech/loggo"
	"github.com/claudetech/loggo/appenders"
	"os"
	"time"

	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	url = kingpin.Arg("URL", "URL address to check.").Required().String()

	reporter  = kingpin.Flag("reporter", "Reporter to use to publish stats (http|redis|stdout)").Default("stdout").Short('r').Enum("stdout", "redis", "http")
	reportUrl = kingpin.Flag("report-url", "URL to report. HTTP(s) when using HTTP or Redis dial info for redis.").Short('u').String()
	queueName = kingpin.Flag("queue-name", "Name of queue to use when using Redis").Short('q').String()

	format = kingpin.Flag("format", "Format to POST stats (only json available for now)").Default("json").Enum("json")

	logLevel  = kingpin.Flag("log-level", "Set log level.").Default("info").Enum("trace", "debug", "info", "warning", "error", "fatal")
	logFile   = kingpin.Flag("log-file", "File to log the output").Short('o').String()
	useStdout = kingpin.Flag("print", "Print the output to stdout").Short('p').Bool()
	slackUrl  = kingpin.Flag("slack-url", "URL to send message to Slack when server is not availabile").String()

	checkInterval   = kingpin.Flag("check-interval", "Interval in seconds between availability check").Default("10").Int()
	publishInterval = kingpin.Flag("publish-interval", "Interval in seconds between stats publication").Default("60").Int()
	fatalThreshold  = kingpin.Flag("fatal-threshold", "Sucess ratio under which a fatal error should be logged").Default("0.8").Short('f').Float()
	extraFields     = kingpin.Flag("extra-fields", "Extra fields to send with the report stats data").Short('x').StringMap()
	requestHeaders  = kingpin.Flag("request-headers", "Extra headers to send with the request").Short('H').StringMap()
	requestType     = kingpin.Flag("request-type", "The method to use when sending message to the server to check").Default("GET").String()
	requestBody     = kingpin.Flag("request-body", "The body to send to the server to check if using POST").String()
)

func getMimetype() string {
	switch *format {
	default:
		return "application/json"
	}
}

func getReporter() (avchecker.Reporter, error) {
	switch *reporter {
	case "redis":
		if *queueName == "" {
			kingpin.UsageErrorf("queue-name must be provided when using Redis")
		}
		if reportUrl == nil {
			return nil, fmt.Errorf("'report-url' is required when using redis, try --help")
		}
		return avchecker.NewRedisQueueReporter(*reporter, *reportUrl)
	case "http":
		if reportUrl == nil {
			return nil, fmt.Errorf("'report-url' is required when using http, try --help")
		}
		return avchecker.NewHttpReporter(*reportUrl, getMimetype())
	case "stdout":
		return avchecker.NewStdoutReporter(), nil
	default:
		return nil, fmt.Errorf("unknown type %s", *reporter)
	}
}

func makeLogger() (*loggo.Logger, error) {
	logger := loggo.New("avchecker")
	logger.SetLevel(loggo.LevelFromString(*logLevel))

	if *useStdout {
		logger.AddAppender(loggo.NewStdoutAppender(), loggo.Color)
	}

	if *logFile != "" {
		appender, err := loggo.NewFileAppender(*logFile)
		if err != nil {
			return nil, err
		}
		logger.AddAppender(appender, loggo.EmptyFlag)
	}

	if !*useStdout && *logFile == "" {
		fmt.Println("Warning: Neither print option or a log file has been passed. Output will be lost.")
	}

	if *slackUrl != "" {
		filter := &loggo.MinLogLevelFilter{MinLevel: loggo.Fatal}
		slack := appenders.NewSlackAppender(*slackUrl, "AvailabilityBot", ":warning:", "")
		logger.AddAppenderWithFilter(slack, filter, loggo.Async)
	}
	return logger, nil
}

func main() {
	kingpin.Version("0.1.0")
	kingpin.Parse()

	reporter, err := getReporter()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error has occured: %s", err)
		os.Exit(1)
	}
	logger, err := makeLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error has occured: %s", err)
		os.Exit(1)
	}

	checker := avchecker.New(*url, reporter, &avchecker.Options{
		CheckInterval:   time.Duration(*checkInterval) * time.Second,
		PublishInterval: time.Duration(*publishInterval) * time.Second,
		Logger:          logger,
		FatalThreshold:  *fatalThreshold,
		RequestType:     *requestType,
		RequestBody:     bytes.NewReader([]byte(*requestBody)),
		RequestHeaders:  *requestHeaders,
		ExtraFields:     *extraFields,
	})

	logger.Infof("Starting to monitor availability for %s", *url)

	checker.StartChecking()
}
