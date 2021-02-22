package step_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jayclassless/dotter/step"
)

var _ = Describe("util", func() {
	var testDir string

	BeforeEach(func() {
		testDir = tmpdir()
	})

	AfterEach(func() {
		rmdir(testDir)
	})

	Describe("IsSymLink", func() {
		It("Works", func() {
			ln(filepath.Join(testDir, "foo"), "/tmp")
			fileInfo, err := os.Lstat(filepath.Join(testDir, "foo"))
			Expect(err).To(Succeed())
			Expect(step.IsSymLink(fileInfo)).To(BeTrue())

			writeFile(testDir, "bar", "bar")
			fileInfo, err = os.Stat(filepath.Join(testDir, "bar"))
			Expect(err).To(Succeed())
			Expect(step.IsSymLink(fileInfo)).To(BeFalse())
		})
	})
})
