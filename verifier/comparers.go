package verifier

import (
	"github.com/VerifyTests/Verify.Go/utils"
	"strings"
)

type compare struct {
}

var comparer = compare{}

func (c *compare) Text(filePair FilePair, received string, settings *verifySettings) EqualityResult {
	utils.File.DeleteIfEmpty(filePair.VerifiedPath)
	if !utils.File.Exists(filePair.VerifiedPath) {
		utils.File.WriteText(filePair.ReceivedPath, received)
		return EqualityResult{
			Equality: FileNew,
		}
	}

	verified := strings.Builder{}
	verified.Write(utils.File.ReadText(filePair.VerifiedPath))

	result := compareStrings(filePair.Extension, received, verified.String(), settings)
	if result.IsEqual {
		return EqualityResult{
			Equality: FileEqual,
		}
	}

	utils.File.WriteText(filePair.ReceivedPath, received)
	return EqualityResult{
		Equality: FileNotEqual,
		Message:  result.Message,
	}
}

func compareStrings(extension string, received string, verified string, settings *verifySettings) CompareResult {
	isEqual := verified == received

	if !isEqual {
		comparer, ok := settings.tryGetStringComparer(extension)
		if ok {
			comparer(received, verified)
		}
	}

	return CompareResult{IsEqual: isEqual}
}

func (c *compare) Streams(FilePair, []byte, *verifySettings, bool) EqualityResult {
	panic("not implemented")
}
