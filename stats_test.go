package avchecker

import (
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = g.Describe("stats", func() {
	g.Describe("compute", func() {
		g.It("should compute success_ratio", func() {
			s := &stats{1, 1, 0, nil}
			s.compute()
			Expect(s.SuccessRatio).To(Equal(1.0))
			s = &stats{2, 1, 0, nil}
			s.compute()
			Expect(s.SuccessRatio).To(Equal(0.5))
		})
	})

	g.Describe("toMap", func() {
		g.It("should work with base stats", func() {
			s := &stats{2, 1, 0, nil}
			m := s.toMap()
			Expect(m).To(HaveKey("success_count"))
		})

		g.It("should work with extra fields", func() {
			s := &stats{2, 1, 0, map[string]interface{}{"foo": "bar"}}
			m := s.toMap()
			Expect(m).To(HaveKey("foo"))
		})
	})
})
