package observer_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/st3v/arp-watch/observer"
)

var _ = Describe("State", func() {
	Describe(".Load", func() {
		var (
			actual        *observer.State
			stateFilePath string
		)

		BeforeEach(func() {
			actual = new(observer.State)
			stateFilePath = filepath.Join("..", "test", "assets", "state.json")
		})

		It("does not return an error", func() {
			err := actual.Load(stateFilePath)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns the expected state", func() {
			actual.Load(stateFilePath)
			Expect(*actual).To(Equal(observer.State{
				"1.2.3.4": "aa:aa:aa:aa:aa:aa",
				"5.6.7.8": "bb:bb:bb:bb:bb:bb",
			}))
		})

		Context("when stateFilePath does not exist", func() {
			BeforeEach(func() {
				stateFilePath = filepath.Join("path", "to", "nowhere")
			})

			It("returns a corresponding error", func() {
				err := actual.Load(stateFilePath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
		})

		Context("when stateFilePath points to an invalid file", func() {
			BeforeEach(func() {
				stateFilePath = filepath.Join("..", "test", "assets", "invalid.json")
			})

			It("returns a corresponding error", func() {
				err := actual.Load(stateFilePath)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid character"))
			})
		})
	})

	Describe(".Write", func() {
		var (
			expected = &observer.State{
				"1.2.3.4": "aa:aa:aa:aa:aa:aa",
				"5.6.7.8": "bb:bb:bb:bb:bb:bb",
			}

			stateFilePath string
		)

		BeforeEach(func() {
			tmpDirPath, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			stateFilePath = filepath.Join(tmpDirPath, "state.json")
		})

		AfterEach(func() {
			os.RemoveAll(filepath.Dir(stateFilePath))
		})

		It("does not return an error", func() {
			err := expected.Write(stateFilePath)
			Expect(err).NotTo(HaveOccurred())
		})

		It("correctly writes the file", func() {
			expected.Write(stateFilePath)

			actual := new(observer.State)
			err := actual.Load(stateFilePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(*actual).To(Equal(*expected))
		})
	})
})
