package verifier

import (
	"github.com/VerifyTests/Verify.Go/utils"
	"strings"
)

type Target struct {
	stringData           string
	stringBuilderData    *strings.Builder
	streamData           []byte
	extension            string
	hasStringBuilderData bool
	hasStringData        bool
	hasStreamData        bool
}

func newStringTarget(extension string, stringData string) *Target {
	utils.Guard.AgainstBadExtension(extension)
	if !utils.File.IsText(extension) {
		panic("Dont pass a text for a binary extension. Instead use `newStreamTarget`")
	}

	return &Target{
		extension:     extension,
		stringData:    stringData,
		hasStringData: true,
	}
}

func newStringBuilderTarget(extension string, stringBuilderData *strings.Builder) *Target {
	utils.Guard.AgainstBadExtension(extension)
	if !utils.File.IsText(extension) {
		panic("Dont pass a text for a binary extension. Instead use `newStreamTarget`")
	}

	return &Target{
		extension:            extension,
		stringBuilderData:    stringBuilderData,
		hasStringBuilderData: true,
	}
}

func newStreamTarget(extension string, stream []byte) *Target {
	utils.Guard.AgainstBadExtension(extension)
	if utils.File.IsText(extension) {
		panic("Dont pass a byte slice for text. Instead use `newStringTarget` or `newStringBuilderTarget`")
	}

	return &Target{
		extension:     extension,
		streamData:    stream,
		hasStreamData: true,
	}
}

func (t *Target) IsStringBuilder() bool {
	return t.hasStringBuilderData
}

func (t *Target) IsString() bool {
	return t.hasStringData
}

func (t *Target) IsStream() bool {
	return t.hasStreamData
}

func (t *Target) GetStringBuilderData() *strings.Builder {
	if !t.hasStringBuilderData {
		panic("Use `GetStreamData` or `GetStringData`")
	}
	return t.stringBuilderData
}

func (t *Target) GetStringData() string {
	if !t.hasStringData {
		panic("Use `GetStreamData` or `GetStringBuilderData`")
	}
	return t.stringData
}

func (t *Target) GetStreamData() []byte {
	if !t.hasStreamData {
		panic("Use `GetStringData` or `GetStringBuilderData`")
	}
	return t.streamData
}

func (t *Target) GetExtension() string {
	return t.extension
}
