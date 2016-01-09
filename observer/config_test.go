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
			Expect(actual.Metron.Endpoint).To(Equal("localhost:3457"))
			Expect(actual.Metron.Origin).To(Equal("node-1"))
			Expect(actual.Frequency).To(Equal("1s"))
			Expect(actual.Filters).To(Equal([]string{"192.168.0.1", "192.168.0.2"}))
			Expect(actual.Aliases).To(Equal(map[string]string{
				"192.168.0.1": "host-1",
				"192.168.0.2": "host-2",
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
