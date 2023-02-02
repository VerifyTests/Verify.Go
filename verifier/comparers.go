package verifier

import (
	"bytes"
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
	verified.Write(utils.File.ReadFile(filePair.VerifiedPath))

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

func (c *compare) Streams(filePair FilePair, receivedStream []byte, settings *verifySettings) EqualityResult {

	if !utils.File.Exists(filePair.VerifiedPath) {
		utils.File.WriteStream(filePair.ReceivedPath, receivedStream)
		return EqualityResult{
			Equality: FileNew,
		}
	}

	if utils.File.GetLength(filePair.VerifiedPath) == 0 {
		utils.File.WriteStream(filePair.ReceivedPath, receivedStream)
		return EqualityResult{
			Equality: FileNotEqual,
		}
	}

	if utils.File.GetLength(filePair.VerifiedPath) != int64(len(receivedStream)) {
		utils.File.WriteStream(filePair.ReceivedPath, receivedStream)
		return EqualityResult{
			Equality: FileNotEqual,
		}
	}

	verifiedStream := utils.File.ReadFile(filePair.VerifiedPath)
	if !bytes.Equal(verifiedStream, receivedStream) {
		utils.File.WriteStream(filePair.ReceivedPath, receivedStream)
		return EqualityResult{
			Equality: FileNotEqual,
		}
	}

	return EqualityResult{
		Equality: FileEqual,
	}
}
