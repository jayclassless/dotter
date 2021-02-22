package step_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jayclassless/dotter/step"
)

var _ = Describe("LinkStep", func() {
	Describe("NewLinkStep", func() {
		It("Works", func() {
			Expect(step.NewLinkStep()).ShouldNot(BeNil())
		})
	})

	Describe("GetActivityLabel", func() {
		It("Works", func() {
			Expect(step.NewLinkStep().GetActivityLabel()).To(Equal("Linking"))
		})
	})

	Describe("GetActivityDetails", func() {
		It("Works", func() {
			step := step.NewLinkStep()
			step.Target = "foobar"

			Expect(step.GetActivityDetails()).To(Equal("foobar"))
		})
	})

	Describe("Execute", func() {
		var executor *TestExecutor

		BeforeEach(func() {
			executor = NewTestExecutor(tmpdir(), tmpdir())
		})

		AfterEach(func() {
			rmdir(executor.target)
			rmdir(executor.source)
		})

		It("Handles the simple case", func() {
			s := step.NewLinkStep()
			s.Target = "foo"
			s.Source = "bar"

			err := s.Execute(executor)
			Expect(err).Should(Succeed())

			fileInfo, err := os.Lstat(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(step.IsSymLink(fileInfo)).To(BeTrue())

			linkPath, err := os.Readlink(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			link := "../" + filepath.Base(executor.source) + "/bar"
			Expect(linkPath).To(Equal(link))
		})

		It("Handles non-relative link", func() {
			s := step.NewLinkStep()
			s.Target = "foo"
			s.Source = "bar"
			s.Relative = false

			err := s.Execute(executor)
			Expect(err).Should(Succeed())

			fileInfo, err := os.Lstat(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(step.IsSymLink(fileInfo)).To(BeTrue())

			linkPath, err := os.Readlink(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(linkPath).To(Equal(executor.GetSourcePath("bar")))
		})

		It("Handles deep paths", func() {
			s := step.NewLinkStep()
			s.Target = "some/deep/foo"
			s.Source = "bar"
			s.Relative = false

			err := s.Execute(executor)
			Expect(err).Should(Succeed())

			fileInfo, err := os.Lstat(executor.GetTargetPath("some/deep/foo"))
			Expect(err).Should(Succeed())
			Expect(step.IsSymLink(fileInfo)).To(BeTrue())

			linkPath, err := os.Readlink(executor.GetTargetPath("some/deep/foo"))
			Expect(err).Should(Succeed())
			Expect(linkPath).To(Equal(executor.GetSourcePath("bar")))
		})

		It("Fails on deep paths when CreateParents is disabled", func() {
			s := step.NewLinkStep()
			s.Target = "some/deep/foo"
			s.Source = "bar"
			s.CreateParents = false

			err := s.Execute(executor)
			Expect(err).Should(HaveOccurred())

			_, err = os.Lstat(executor.GetTargetPath("some/deep/foo"))
			Expect(err).Should(HaveOccurred())
		})

		It("Handles existing links", func() {
			s := step.NewLinkStep()
			s.Target = "foo"
			s.Source = "bar"
			s.Relative = false

			ln(executor.GetTargetPath("foo"), executor.GetSourcePath("bar"))

			err := s.Execute(executor)
			Expect(err).Should(Succeed())

			fileInfo, err := os.Lstat(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(step.IsSymLink(fileInfo)).To(BeTrue())

			linkPath, err := os.Readlink(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(linkPath).To(Equal(executor.GetSourcePath("bar")))
		})

		It("Relinks existing links", func() {
			s := step.NewLinkStep()
			s.Target = "foo"
			s.Source = "bar"
			s.Relative = false

			ln(executor.GetTargetPath("foo"), "bogus")

			err := s.Execute(executor)
			Expect(err).Should(Succeed())

			fileInfo, err := os.Lstat(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(step.IsSymLink(fileInfo)).To(BeTrue())

			linkPath, err := os.Readlink(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(linkPath).To(Equal(executor.GetSourcePath("bar")))
		})

		It("Fails on existing links when Relink is disabled", func() {
			s := step.NewLinkStep()
			s.Target = "foo"
			s.Source = "bar"
			s.Relink = false

			ln(executor.GetTargetPath("foo"), "bogus")

			err := s.Execute(executor)
			Expect(err).Should(HaveOccurred())

			fileInfo, err := os.Lstat(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(step.IsSymLink(fileInfo)).To(BeTrue())

			linkPath, err := os.Readlink(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(linkPath).To(Equal("bogus"))
		})

		It("Handles collisions when Force is enabled", func() {
			s := step.NewLinkStep()
			s.Target = "foo"
			s.Source = "bar"
			s.Relative = false
			s.Force = true

			writeFile(executor.target, "foo", "foo")

			err := s.Execute(executor)
			Expect(err).Should(Succeed())
			Expect(executor.backedUp).To(HaveLen(1))

			fileInfo, err := os.Lstat(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(step.IsSymLink(fileInfo)).To(BeTrue())

			linkPath, err := os.Readlink(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(linkPath).To(Equal(executor.GetSourcePath("bar")))
		})

		It("Fails on collisions when Force is disabled", func() {
			s := step.NewLinkStep()
			s.Target = "foo"
			s.Source = "bar"
			s.Relative = false
			s.Force = false

			writeFile(executor.target, "foo", "foo")

			err := s.Execute(executor)
			Expect(err).Should(HaveOccurred())
			Expect(executor.backedUp).To(HaveLen(0))

			fileInfo, err := os.Stat(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(fileInfo.Mode().IsRegular()).To(BeTrue())
		})
	})
})
