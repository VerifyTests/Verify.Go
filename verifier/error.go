package verifier

import (
	"fmt"
	"github.com/VerifyTests/Verify.Go/utils"
	"io"
	"strings"
)

type failingMessageBuilder struct {
	settings      *verifySettings
	directory     string
	notEqualFiles NotEqualFiles
	equalFiles    EqualFiles
	newFiles      NewFiles
	delete        []string
	testCase      string
	testName      string
}

func (b *failingMessageBuilder) build() string {
	builder := strings.Builder{}

	builder.WriteRune('\n')

	if len(b.testName) > 0 {
		builder.WriteString(fmt.Sprintf("Test Name: %s\n", b.testName))
	}

	if len(b.testCase) > 0 {
		builder.WriteString(fmt.Sprintf("Test Case: %s\n", b.testCase))
	}

	builder.WriteString(fmt.Sprintf("Directory: %s\n", b.directory))

	if len(b.newFiles) > 0 {
		builder.WriteString("New:\n")
		for _, f := range b.newFiles {
			b.appendFile(&builder, f)
		}
	}

	if len(b.notEqualFiles) > 0 {
		builder.WriteString("NotEqual:\n")
		for _, f := range b.notEqualFiles {
			b.appendFile(&builder, f.File)
		}
	}

	if len(b.delete) > 0 {
		builder.WriteString("Delete:\n")
		for _, f := range b.delete {
			builder.WriteString(fmt.Sprintf("  - %s", utils.File.GetFileName(f)))
		}
	}

	if len(b.equalFiles) > 0 {
		builder.WriteString("Equal:\n")
		for _, f := range b.equalFiles {
			b.appendFile(&builder, f)
		}
	}

	b.appendContent(&builder)

	return builder.String()
}

func (b *failingMessageBuilder) appendFile(builder io.StringWriter, file FilePair) {
	_, _ = builder.WriteString(fmt.Sprintf("  - Received: %s\n", file.ReceivedName))
	_, _ = builder.WriteString(fmt.Sprintf("    Verified: %s\n", file.VerifiedName))
}

func (b *failingMessageBuilder) appendContent(builder *strings.Builder) {
	if b.settings.omitContentFromError {
		return
	}

	newContentFiles := make([]FilePair, 0)
	notEqualContentFiles := make([]NotEqualFilePair, 0)

	for _, f := range b.newFiles {
		if f.IsText {
			newContentFiles = append(newContentFiles, f)
		}
	}

	for _, f := range b.notEqualFiles {
		if f.File.IsText || len(f.Message) > 0 {
			notEqualContentFiles = append(notEqualContentFiles, f)
		}
	}

	if len(newContentFiles) == 0 && len(notEqualContentFiles) == 0 {
		return
	}

	builder.WriteString("\nFileContent:\n")

	if len(newContentFiles) > 0 {
		builder.WriteString("New:\n")
		for _, item := range newContentFiles {
			builder.WriteString(fmt.Sprintf("  - Received: %s\n", item.ReceivedName))
			builder.Write(utils.File.ReadFile(item.ReceivedPath))
			builder.WriteString("\n")
		}
	}

	if len(notEqualContentFiles) > 0 {
		builder.WriteString("  NotEqual:\n")
		for _, item := range notEqualContentFiles {
			if len(item.Message) == 0 {
				builder.WriteString(fmt.Sprintf("  - Received: %s\n", item.File.ReceivedName))
				builder.Write(utils.File.ReadFile(item.File.ReceivedPath))
				builder.WriteRune('\n')
				builder.WriteString(fmt.Sprintf("    Verified: %s\n", item.File.ReceivedName))
				builder.Write(utils.File.ReadFile(item.File.VerifiedPath))
				builder.WriteRune('\n')
			} else {
				builder.WriteString(fmt.Sprintf("  - Received: %s\n", item.File.ReceivedName))
				builder.WriteString(fmt.Sprintf("    Verified: %s\n", item.File.VerifiedName))
				builder.WriteString(fmt.Sprintf("    Compare Result: %s\n", item.Message))
			}
			builder.WriteString("\n")
		}
	}
}
