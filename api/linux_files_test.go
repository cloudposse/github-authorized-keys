package api

import (
	"github.com/goruha/permbits"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Linux files", func() {

	var linux Linux

	BeforeEach(func() {
		linux = NewLinux("/")
	})

	Describe("FileEnsure()", func() {
		Context("call with non-existing file", func() {

			AfterEach(func() {
				linux.FileDelete("/tmp/zzz")
			})

			It("should create file with right content", func() {
				isFound := linux.FileExists("/tmp/zzz")

				Expect(isFound).To(BeFalse())

				linux.FileEnsure("/tmp/zzz", "RIGHT CONTENT")

				isFound = linux.FileExists("/tmp/zzz")

				Expect(isFound).To(BeTrue())

				content, _ := linux.FileGet("/tmp/zzz")

				Expect(content).To(Equal("RIGHT CONTENT"))

			})
		})
	})

	Describe("FileExists()", func() {
		Context("call with existing file", func() {
			It("should return true", func() {
				isFound := linux.FileExists("/bin/sh")

				Expect(isFound).To(BeTrue())
			})
		})
	})

	Describe("FileCreate()", func() {
		Context("call with non-existing file", func() {

			AfterEach(func() {
				linux.FileDelete("/tmp/zzz")
			})

			It("should create file", func() {
				isFound := linux.FileExists("/tmp/zzz")

				Expect(isFound).To(BeFalse())

				linux.FileCreate("/tmp/zzz")

				isFound = linux.FileExists("/tmp/zzz")

				Expect(isFound).To(BeTrue())
			})
		})
	})

	Describe("FileGet()", func() {
		Context("call with existing non-empty file", func() {
			It("should return no-empty content", func() {
				content, err := linux.FileGet("/bin/sh")

				Expect(err).To(BeNil())
				Expect(content).NotTo(BeEmpty())
			})
		})
	})

	Describe("FileSet()", func() {
		Context("call with non-existing file", func() {

			AfterEach(func() {
				linux.FileDelete("/tmp/zzz")
			})

			It("should create file with valid content", func() {
				err := linux.FileSet("/tmp/zzz", "RIGHT CONTENT")
				Expect(err).To(BeNil())

				content, err := linux.FileGet("/tmp/zzz")

				Expect(err).To(BeNil())
				Expect(content).To(Equal("RIGHT CONTENT"))
			})
		})
	})

	Describe("FileModeGet()", func() {
		Context("call with existing file with x perm", func() {
			It("should get all executable perm", func() {
				mode, err := linux.FileModeGet("/bin/sh")

				Expect(err).To(BeNil())
				Expect(mode.OtherExecute()).To(BeTrue())
				Expect(mode.GroupExecute()).To(BeTrue())
				Expect(mode.UserExecute()).To(BeTrue())

			})
		})
	})

	Describe("FileModeSet()", func() {
		Context("call with existing file with wrong permissions", func() {

			BeforeEach(func() {
				linux.FileCreate("/tmp/zzz")
			})

			AfterEach(func() {
				linux.FileDelete("/tmp/zzz")
			})

			It("should set right perms", func() {
				mode, err := linux.FileModeGet("/tmp/zzz")

				Expect(err).To(BeNil())
				Expect(mode.OtherExecute()).To(BeFalse())

				mode.SetOtherExecute(true)
				err = linux.FileModeSet("/tmp/zzz", mode)
				Expect(err).To(BeNil())

				mode, err = linux.FileModeGet("/tmp/zzz")
				Expect(err).To(BeNil())
				Expect(mode.OtherExecute()).To(BeTrue())

			})
		})
	})

	Describe("FileModeEnsure()", func() {
		Context("call with existing file with x perm", func() {
			BeforeEach(func() {
				linux.FileCreate("/tmp/zzz")
				mode, _ := linux.FileModeGet("/tmp/zzz")
				mode.SetOtherExecute(true)
				linux.FileModeSet("/tmp/zzz", mode)

			})

			AfterEach(func() {
				linux.FileDelete("/tmp/zzz")
			})

			It("should set right perms", func() {
				mode, err := linux.FileModeGet("/tmp/zzz")

				Expect(err).To(BeNil())

				Expect(mode.UserWrite()).To(BeTrue())
				Expect(mode.UserRead()).To(BeTrue())
				Expect(mode.UserExecute()).To(BeFalse())

				Expect(mode.GroupWrite()).To(BeFalse())
				Expect(mode.GroupRead()).To(BeTrue())
				Expect(mode.GroupExecute()).To(BeFalse())

				Expect(mode.OtherWrite()).To(BeFalse())
				Expect(mode.OtherRead()).To(BeTrue())
				Expect(mode.OtherExecute()).To(BeTrue())

				requiredMode := permbits.PermissionBits(0)
				requiredMode.SetGroupExecute(true)
				err = linux.FileModeEnsure("/tmp/zzz", requiredMode)
				Expect(err).To(BeNil())

				mode, err = linux.FileModeGet("/tmp/zzz")
				Expect(err).To(BeNil())

				Expect(mode.UserWrite()).To(BeTrue())
				Expect(mode.UserRead()).To(BeTrue())
				Expect(mode.UserExecute()).To(BeFalse())

				Expect(mode.GroupWrite()).To(BeFalse())
				Expect(mode.GroupRead()).To(BeTrue())
				Expect(mode.GroupExecute()).To(BeTrue())

				Expect(mode.OtherWrite()).To(BeFalse())
				Expect(mode.OtherRead()).To(BeTrue())
				Expect(mode.OtherExecute()).To(BeTrue())
			})
		})
	})

	Describe("FileEnsureLine()", func() {
		Context("call with file that does not contain target string", func() {
			BeforeEach(func() {
				linux.FileEnsure("/tmp/zzz", "RIGHT22 CONTENT")
			})

			AfterEach(func() {
				linux.FileDelete("/tmp/zzz")
			})

			It("should add target string into file", func() {
				err := linux.FileEnsureLine("/tmp/zzz", "RIGHT CONTENT")
				content, _ := linux.FileGet("/tmp/zzz")

				Expect(err).To(BeNil())
				Expect(content).To(Equal(
					`RIGHT22 CONTENT
RIGHT CONTENT`))
			})
		})

		Context("call with file that contains target string", func() {
			BeforeEach(func() {
				linux.FileEnsure("/tmp/zzz", "RIGHT CONTENT")
			})

			AfterEach(func() {
				linux.FileDelete("/tmp/zzz")
			})

			It("should do nothing with file", func() {
				err := linux.FileEnsureLine("/tmp/zzz", "RIGHT CONTENT")
				content, _ := linux.FileGet("/tmp/zzz")

				Expect(err).To(BeNil())
				Expect(content).To(Equal("RIGHT CONTENT"))
			})
		})

	})

	Describe("FileEnsureLineMatcher()", func() {
		Context("call with file that does not contain target string", func() {
			BeforeEach(func() {
				linux.FileEnsure("/tmp/zzz", "RIGHT22 CONTENT")
			})

			AfterEach(func() {
				linux.FileDelete("/tmp/zzz")
			})

			It("should add target string into file", func() {
				err := linux.FileEnsureLineMatch("/tmp/zzz", "RIGHT\\s.*", "RIGHT CONTENT")
				content, _ := linux.FileGet("/tmp/zzz")

				Expect(err).To(BeNil())
				Expect(content).To(Equal(
					`RIGHT22 CONTENT
RIGHT CONTENT`))
			})
		})

		Context("call with file that contains target string", func() {
			BeforeEach(func() {
				linux.FileEnsure("/tmp/zzz", "RIGHT CONTENT")
			})

			AfterEach(func() {
				linux.FileDelete("/tmp/zzz")
			})

			It("should do nothing with file", func() {
				err := linux.FileEnsureLineMatch("/tmp/zzz", "RIGHT\\s.*", "RIGHT CONTENT")
				content, _ := linux.FileGet("/tmp/zzz")

				Expect(err).To(BeNil())
				Expect(content).To(Equal("RIGHT CONTENT"))
			})
		})


		Context("call with file that contains string that satisfies match but differs from traget", func() {
			BeforeEach(func() {
				linux.FileEnsure("/tmp/zzz", "RIGHT CONTENT22")
			})

			AfterEach(func() {
				linux.FileDelete("/tmp/zzz")
			})

			It("should do replace matched string with target string", func() {
				err := linux.FileEnsureLineMatch("/tmp/zzz", "RIGHT\\s.*", "RIGHT CONTENT")
				content, _ := linux.FileGet("/tmp/zzz")

				Expect(err).To(BeNil())
				Expect(content).To(Equal("RIGHT CONTENT"))
			})
		})

		Context("call with too common matcher", func() {
			BeforeEach(func() {
				linux.FileEnsure("/tmp/zzz", "RIGHT CONTENT23")
				linux.FileEnsureLine("/tmp/zzz", "RIGHT CONTENT")
			})

			AfterEach(func() {
				linux.FileDelete("/tmp/zzz")
			})

			It("should retrun valid error", func() {
				err := linux.FileEnsureLineMatch("/tmp/zzz", ".*", "RIGHT CONTENT")
				linux.FileGet("/tmp/zzz")

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("Match regexp /.*/ is too wide - [RIGHT CONTENT23 RIGHT CONTENT] matches found."))
			})
		})


	})

})
