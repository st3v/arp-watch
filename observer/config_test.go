package observer_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/st3v/arp-watch/observer"
)

var _ = Describe("Config", func() {

	Describe("load", func() {
		var (
			actual     *observer.Config
			configPath string
		)

		BeforeEach(func() {
			actual = new(observer.Config)
			configPath = filepath.Join("..", "test", "assets", "config.json")
		})

		It("does not return an error", func() {
			err := actual.Load(configPath)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns the expected config", func() {
			actual.Load(configPath)
			Expect(actual.Metron.Endpoint).To(Equal("endpoint"))
			Expect(actual.Metron.Origin).To(Equal("origin"))
			Expect(actual.Frequency).To(Equal("1s"))
			Expect(actual.Filter).To(Equal([]string{"1.2.3.4", "5.6.7.8"}))
			Expect(actual.Alias).To(Equal(map[string]string{
				"1.2.3.4": "one",
				"5.6.7.8": "two",
			}))
		})

		Context("when configPath does not exist", func() {
			BeforeEach(func() {
				configPath = filepath.Join("path", "to", "nowhere")
			})

			It("returns a corresponding error", func() {
				err := actual.Load(configPath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
		})

		Context("when configPath points to an invalid file", func() {
			BeforeEach(func() {
				configPath = filepath.Join("..", "test", "assets", "invalid.json")
			})

			It("returns a corresponding error", func() {
				err := actual.Load(configPath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid character"))
			})
		})
	})
})
