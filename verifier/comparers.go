package verifier

import "strings"

type compare struct {
}

var comparer = compare{}

func (c *compare) Text(filePair FilePair, received string, settings *verifySettings) EqualityResult {
	file.deleteIfEmpty(filePair.VerifiedPath)
	if !file.exists(filePair.VerifiedPath) {
		file.writeText(filePair.ReceivedPath, received)
		return EqualityResult{
			Equality: FileNew,
		}
	}

	verified := strings.Builder{}
	verified.Write(file.readText(filePair.VerifiedPath))

	result := compareStrings(filePair.Extension, received, verified.String(), settings)
	if result.IsEqual {
		return EqualityResult{
			Equality: FileEqual,
		}
	}

	file.writeText(filePair.ReceivedPath, received)
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
