package verifier

import (
	"github.com/VerifyTests/Verify.Go/diff"
	"github.com/VerifyTests/Verify.Go/utils"
	"strings"
)

type engine struct {
	directory           string
	deletedFiles        []string
	testing             testingT
	settings            *verifySettings
	getIndexedFileNames getIndexedFileNamesFunc
	getFileNames        getFileNamesFunc
	equalFiles          EqualFiles
	newFiles            NewFiles
	notEqualFiles       NotEqualFiles
	detector            diff.CIDetected
}

func newEngine(
	testing testingT,
	directory string,
	settings *verifySettings,
	verifiedFiles []string,
	getFileNames getFileNamesFunc,
	getIndexedFileNames getIndexedFileNamesFunc) *engine {

	return &engine{
		testing:             testing,
		directory:           directory,
		settings:            settings,
		getFileNames:        getFileNames,
		getIndexedFileNames: getIndexedFileNames,
		deletedFiles:        verifiedFiles,
		notEqualFiles:       NotEqualFiles{},
		equalFiles:          EqualFiles{},
		newFiles:            NewFiles{},
	}
}

func (e *engine) handleResults(targetList []Target) {
	if len(targetList) == 1 {
		target := targetList[0]
		file := e.getFileNames(target.GetExtension())
		result := e.getResult(file, target, false)

		e.handleCompareResult(result, file)
		return
	}

	textHasFailed := false
	for index := 0; index < len(targetList); index++ {
		target := targetList[index]
		file := e.getIndexedFileNames(target.GetExtension(), index)
		result := e.getResult(file, target, textHasFailed)

		if file.IsText && result.Equality != FileEqual {
			textHasFailed = true
		}

		e.handleCompareResult(result, file)
	}
}

func (e *engine) throwIfRequired() {
	e.processEquals()

	noChanges := len(e.newFiles) == 0 &&
		len(e.notEqualFiles) == 0 &&
		len(e.deletedFiles) == 0

	if noChanges {
		return
	}

	e.processDeletes()

	e.processNew()

	e.processNotEquals()

	if !e.settings.autoVerify {
		errorBuilder := e.newErrorBuilder()
		e.testing.Errorf(errorBuilder.build())
	}
}

func (e *engine) newErrorBuilder() failingMessageBuilder {
	return failingMessageBuilder{
		settings:      e.settings,
		testCase:      e.settings.testCase,
		testName:      e.testing.Name(),
		directory:     e.directory,
		notEqualFiles: e.notEqualFiles,
		equalFiles:    e.equalFiles,
		newFiles:      e.newFiles,
		delete:        e.deletedFiles,
	}
}

func (e *engine) handleCompareResult(compareResult EqualityResult, file FilePair) {
	switch compareResult.Equality {
	case FileNew:
		e.addMissing(file)
	case FileNotEqual:
		e.addNotEquals(file, compareResult.Message)
	case FileEqual:
		e.addEquals(file)
	}
}

func (e *engine) addMissing(item FilePair) {
	e.newFiles = append(e.newFiles, item)
	e.deletedFiles = removeStringItem(e.deletedFiles, item.VerifiedPath)
}

func (e *engine) addNotEquals(item FilePair, message string) {
	neq := NotEqualFilePair{
		File:    item,
		Message: message,
	}
	e.notEqualFiles = append(e.notEqualFiles, neq)
}

func (e *engine) addEquals(item FilePair) {
	e.deletedFiles = removeStringItem(e.deletedFiles, item.VerifiedPath)
	e.equalFiles = append(e.equalFiles, item)
}

func (e *engine) processDeletes() {
	if len(e.deletedFiles) == 0 {
		return
	}

	for _, item := range e.deletedFiles {
		e.processDelete(item)
	}
}

func (e *engine) processDelete(deletedFile string) {

	e.settings.runOnVerifyDelete(deletedFile)

	if e.settings.autoVerify {
		utils.File.Delete(deletedFile)
		return
	}

	if e.settings.ciDetected {
		return
	}

	//TODO: how to implement diff tray?
}

func (e *engine) processNew() {
	if len(e.newFiles) == 0 {
		return
	}

	for _, item := range e.newFiles {
		e.settings.runOnFirstVerify(item)
		e.runDiffAutoCheck(item)
	}
}

func (e *engine) processNotEquals() {
	if len(e.notEqualFiles) == 0 {
		return
	}

	for _, item := range e.notEqualFiles {
		e.settings.runOnVerifyMismatch(item.File, item.Message)
		e.runDiffAutoCheck(item.File)
	}
}

func (e *engine) processEquals() {
	if !e.settings.diffDisabled {
		return
	}

	for _, item := range e.equalFiles {
		diff.Kill(item.ReceivedPath, item.VerifiedPath)
	}
}

func (e *engine) getResult(filePair FilePair, target Target, previousTextHasFailed bool) EqualityResult {
	if target.IsStringBuilder() {
		builder := target.GetStringBuilderData()
		e.settings.scrubber.Apply(target.GetExtension(), builder, e.settings)
		return comparer.Text(filePair, builder.String(), e.settings)
	}

	if target.IsString() {
		builder := strings.Builder{}
		builder.WriteString(target.GetStringData())
		e.settings.scrubber.Apply(target.GetExtension(), &builder, e.settings)
		return comparer.Text(filePair, builder.String(), e.settings)
	}

	var stream = target.GetStreamData()
	return comparer.Streams(filePair, stream, e.settings, previousTextHasFailed)
}

func (e *engine) runDiffAutoCheck(item FilePair) {
	if e.settings.ciDetected == true {
		return
	}

	if e.settings.autoVerify {
		e.acceptChanges(item)
		return
	}

	if e.settings.diffDisabled {
		diff.Launch(item.ReceivedPath, item.VerifiedPath)
	}
}

func (e *engine) acceptChanges(item FilePair) {
	utils.File.Delete(item.VerifiedPath)
	utils.File.Move(item.ReceivedPath, item.VerifiedPath)
}