package avchecker

import (
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = g.Describe("stats", func() {
	g.Describe("compute", func() {
		g.It("should compute success_ratio", func() {
			s := &stats{1, 1, 0}
			s.compute()
			Expect(s.SuccessRatio).To(Equal(1.0))
			s = &stats{2, 1, 0}
			s.compute()
			Expect(s.SuccessRatio).To(Equal(0.5))
		})
	})
})
