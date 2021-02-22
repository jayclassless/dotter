package dotter_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jayclassless/dotter"
)

var _ = Describe("Config", func() {
	Describe("NewConfiguration", func() {
		It("Works", func() {
			Expect(dotter.NewConfiguration()).ShouldNot(BeNil())
		})
	})

	Describe("NewConfigurationFromFile", func() {
		var tmpDir string

		BeforeEach(func() {
			tmpDir = tmpdir()
		})

		AfterEach(func() {
			rmdir(tmpDir)
			tmpDir = ""
		})

		It("Handles empty files", func() {
			writeFile(tmpDir, "dotter.yaml", "")
			cfg, err := dotter.NewConfigurationFromFile(filepath.Join(tmpDir, "dotter.yaml"))
			Expect(err).Should((Succeed()))

			fresh := dotter.NewConfiguration()
			Expect(cfg.Options).To(Equal(fresh.Options))
			Expect(cfg.Steps).To(Equal(fresh.Steps))
			Expect(cfg.SourcePath).ToNot(Equal(""))
		})

		It("Fails on bad paths", func() {
			_, err := dotter.NewConfigurationFromFile(filepath.Join(tmpDir, "doesntexist.yaml"))
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("NewConfigurationFromYaml", func() {

	})
})
