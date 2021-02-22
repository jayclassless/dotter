package step_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jayclassless/dotter/step"
)

var _ = Describe("DirectoryStep", func() {
	Describe("NewDirectoryStep", func() {
		It("Works", func() {
			Expect(step.NewDirectoryStep()).ShouldNot(BeNil())
		})
	})

	Describe("GetActivityLabel", func() {
		It("Works", func() {
			Expect(step.NewDirectoryStep().GetActivityLabel()).To(Equal("Directory"))
		})
	})

	Describe("GetActivityDetails", func() {
		It("Works", func() {
			step := step.NewDirectoryStep()
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
			step := step.NewDirectoryStep()
			step.Target = "foo"

			err := step.Execute(executor)
			Expect(err).Should((Succeed()))

			fileInfo, err := os.Stat(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(fileInfo.IsDir()).To(BeTrue())
			Expect(fileInfo.Mode().Perm()).To(Equal(os.FileMode(0o755)))
		})

		It("Handles specified mode", func() {
			step := step.NewDirectoryStep()
			step.Target = "foo"
			step.Mode = 0o611

			err := step.Execute(executor)
			Expect(err).Should((Succeed()))

			fileInfo, err := os.Stat(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(fileInfo.IsDir()).To(BeTrue())
			Expect(fileInfo.Mode().Perm()).To(Equal(os.FileMode(0o611)))
		})

		It("Sets specified mode on existing directories", func() {
			step := step.NewDirectoryStep()
			step.Target = "foo"
			step.Mode = 0o611

			mkdir(executor.GetTargetPath("foo"))
			fileInfo, err := os.Stat(executor.GetTargetPath("foo"))
			Expect(fileInfo.Mode().Perm()).ToNot(Equal(os.FileMode(0o611)))

			err = step.Execute(executor)
			Expect(err).Should((Succeed()))

			fileInfo, err = os.Stat(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(fileInfo.IsDir()).To(BeTrue())
			Expect(fileInfo.Mode().Perm()).To(Equal(os.FileMode(0o611)))
		})

		It("Handles deep directories", func() {
			step := step.NewDirectoryStep()
			step.Target = "foo/bar/baz"

			err := step.Execute(executor)
			Expect(err).Should((Succeed()))

			fileInfo, err := os.Stat(executor.GetTargetPath("foo/bar/baz"))
			Expect(err).Should(Succeed())
			Expect(fileInfo.IsDir()).To(BeTrue())
			Expect(fileInfo.Mode().Perm()).To(Equal(os.FileMode(0o755)))
		})

		It("Fails on deep directories when CreateParents is disabled", func() {
			step := step.NewDirectoryStep()
			step.CreateParents = false
			step.Target = "foo/bar/baz"

			err := step.Execute(executor)
			Expect(err).Should(HaveOccurred())

			_, err = os.Stat(executor.GetTargetPath("foo/bar/baz"))
			Expect(os.IsNotExist(err)).To(BeTrue())
		})

		It("Handles collisions when Force is enabled", func() {
			step := step.NewDirectoryStep()
			step.Target = "foo"
			step.Force = true
			writeFile(executor.target, "foo", "foo")

			err := step.Execute(executor)
			Expect(err).Should(Succeed())
			Expect(executor.backedUp).To(HaveLen(1))

			fileInfo, err := os.Stat(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(fileInfo.IsDir()).To(BeTrue())
			Expect(fileInfo.Mode().Perm()).To(Equal(os.FileMode(0o755)))
		})

		It("Fails on collisions when Force is disabled", func() {
			step := step.NewDirectoryStep()
			step.Target = "foo"
			writeFile(executor.target, "foo", "foo")

			err := step.Execute(executor)
			Expect(err).Should(HaveOccurred())
			Expect(executor.backedUp).To(HaveLen(0))

			fileInfo, err := os.Stat(executor.GetTargetPath("foo"))
			Expect(err).Should(Succeed())
			Expect(fileInfo.Mode().IsRegular()).To(BeTrue())
		})
	})
})
