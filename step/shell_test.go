package step_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jayclassless/dotter/step"
)

var _ = Describe("ShellStep", func() {
	Describe("NewShellStep", func() {
		It("Works", func() {
			Expect(step.NewShellStep()).ShouldNot(BeNil())
		})
	})

	Describe("GetActivityLabel", func() {
		It("Works", func() {
			Expect(step.NewShellStep().GetActivityLabel()).To(Equal("Executing"))
		})
	})

	Describe("GetActivityDetails", func() {
		It("Shows command when no description", func() {
			step := step.NewShellStep()
			step.Command = "foobar"

			Expect(step.GetActivityDetails()).To(Equal("foobar"))
		})

		It("Shows description when specified", func() {
			step := step.NewShellStep()
			step.Command = "foobar"
			step.Description = "My Description"

			Expect(step.GetActivityDetails()).To(Equal("My Description"))
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
			step := step.NewShellStep()
			step.Command = "true"

			err := step.Execute(executor)
			Expect(err).Should(Succeed())
			Expect(executor.infoLog).To(HaveLen(0))
			Expect(executor.errorLog).To(HaveLen(0))
		})

		It("Handles a failure", func() {
			step := step.NewShellStep()
			step.Command = "false"

			err := step.Execute(executor)
			Expect(err).Should(HaveOccurred())
			Expect(executor.infoLog).To(HaveLen(0))
			Expect(executor.errorLog).To(HaveLen(0))
		})

		It("Doesn't capture output", func() {
			step := step.NewShellStep()
			step.Command = "echo \"foo\""

			err := step.Execute(executor)
			Expect(err).Should(Succeed())
			Expect(executor.infoLog).To(HaveLen(0))
			Expect(executor.errorLog).To(HaveLen(0))
		})

		It("Captures output when configured", func() {
			step := step.NewShellStep()
			step.Command = "echo \"foo\""
			step.Quiet = false

			err := step.Execute(executor)
			Expect(err).Should(Succeed())
			Expect(executor.infoLog).To(Equal([]string{"foo\n"}))
			Expect(executor.errorLog).To(HaveLen(0))
		})

		It("Captures error output", func() {
			step := step.NewShellStep()
			step.Command = "echo \"foo\" && echo \"bar\" >&2"
			step.Quiet = false

			err := step.Execute(executor)
			Expect(err).Should(Succeed())
			Expect(executor.infoLog).To(Equal([]string{"foo\n"}))
			Expect(executor.errorLog).To(Equal([]string{"bar\n"}))
		})
	})
})
