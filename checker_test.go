package avchecker

import (
	"encoding/json"
	"fmt"
	"github.com/claudetech/loggo"
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"time"
)

type successClient struct{}

func (s *successClient) Do(request *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
	}, nil
}

type failClient struct{}

func (s *failClient) Do(request *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
	}, nil
}

type errorClient struct{}

func (s *errorClient) Do(request *http.Request) (*http.Response, error) {
	return &http.Response{}, fmt.Errorf("error")
}

type dummyReporter struct {
	reports []map[string]interface{}
}

func (r *dummyReporter) SendStats(stats []byte) error {
	var report map[string]interface{}
	if err := json.Unmarshal(stats, &report); err != nil {
		return err
	}
	r.reports = append(r.reports, report)
	return nil
}

func (d *dummyReporter) String() string {
	return "dummy"
}

func checkReports(reports []map[string]interface{}, expected [][]int) {
	Expect(len(reports)).To(Equal(len(expected)))
	for i, expectedVals := range expected {
		report := reports[i]
		Expect(int(report["try_count"].(float64))).To(Equal(expectedVals[0]))
		Expect(int(report["success_count"].(float64))).To(Equal(expectedVals[1]))
	}
}

type strAppender struct {
	s string
}

func (s *strAppender) Append(msg *loggo.Message) {
	s.s += msg.String()
}

var _ = g.Describe("Checker", func() {
	g.Describe("StartChecking", func() {
		var reporter *dummyReporter
		var appender *strAppender

		var makeDummyChecker = func(client miniHttpClient, reporter Reporter, runs int, wait time.Duration) *Checker {
			logger := loggo.New("logger")
			logger.AddAppender(appender, loggo.EmptyFlag)
			return NewChecker("foo", reporter, &Options{
				HttpClient:      client,
				CheckInterval:   1 * time.Nanosecond,
				PublishInterval: wait,
				Logger:          logger,
				totalRuns:       runs,
				ExtraFields:     map[string]string{"foo": "bar"},
				FatalThreshold:  0.5,
			})
		}

		g.BeforeEach(func() {
			reporter = &dummyReporter{}
			appender = &strAppender{}
		})

		g.It("should register successes", func() {
			checker := makeDummyChecker(&successClient{}, reporter, 1, 1*time.Nanosecond)
			checker.StartChecking()
			checkReports(reporter.reports, [][]int{[]int{1, 1}})
			Expect(reporter.reports[0]["success_ratio"].(float64)).To(BeNumerically(">", 0.0))
			Expect(reporter.reports[0]).To(HaveKey("foo"))
		})

		g.It("should work on server failures", func() {
			checker := makeDummyChecker(&failClient{}, reporter, 1, 1*time.Nanosecond)
			checker.StartChecking()
			checkReports(reporter.reports, [][]int{[]int{1, 0}})
		})

		g.It("should work on errors", func() {
			checker := makeDummyChecker(&errorClient{}, reporter, 1, 1*time.Nanosecond)
			checker.StartChecking()
			checkReports(reporter.reports, [][]int{[]int{1, 0}})
		})

		g.It("should send reports only when time is ellapsed", func() {
			checker := makeDummyChecker(&successClient{}, reporter, 3, 10*time.Millisecond)
			checker.StartChecking()
			expected := [][]int{{3, 3}}
			checkReports(reporter.reports, expected)
		})

		g.It("should log fatal when under threshold", func() {
			checker := makeDummyChecker(&failClient{}, reporter, 1, 1*time.Nanosecond)
			Expect(appender.s).To(BeEmpty())
			checker.StartChecking()
			Expect(appender.s).NotTo(BeEmpty())
			Expect(appender.s).To(ContainSubstring("FATAL"))
		})
	})
})
