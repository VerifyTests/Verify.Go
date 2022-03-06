package verifier

import (
	"reflect"
)

// Verifier is the main interface for verification process.
type Verifier interface {
	Verify(target interface{})
}

// Verify verifies the passed target with the default settings.
func Verify(t testingT, target interface{}) {
	VerifyWithSetting(t, NewSettings(), target)
}

// VerifyWithSetting verifies the passed target with the provided settings.
func VerifyWithSetting(t testingT, settings VerifySettings, target interface{}) {
	verifier := NewVerifier(t, settings)
	verifier.Verify(target)
}

// Verify verifies the passed target with the associated settings
func (v *verifier) Verify(target interface{}) {
	defer v.settings.runAfterVerify()

	if isNil(target) {
		v.assertExtensionIsNull()
		v.inner.verifyInner("nil", nil, emptyTargets)
		return
	}

	if stringResult, ok := v.inner.tryGetToString(target); ok {
		if len(stringResult.Extension) > 0 {
			v.settings.UseExtension(stringResult.Extension)
		}

		value := stringResult.Value
		if len(value) == 0 {
			v.inner.verifyInner("emptyString", nil, emptyTargets)
			return
		}

		v.inner.verifyInner(value, nil, emptyTargets)
		return
	}

	if target, ok := target.([]byte); ok {
		v.inner.verifyStream(target)
		return
	}

	v.assertExtensionIsNull()
	v.inner.verifyInner(target, nil, emptyTargets)
}

func (v *verifier) assertExtensionIsNull() {
	if v.settings.extension == "" {
		return
	}

	panic("`UseExtension` should only be used for text, for streams, or for converter discovery.\nWhen serializing an instance the default is txt.\nTo use json as an extension when serializing use `UseStrictJSON`")
}

type verifier struct {
	settings *verifySettings
	inner    *innerVerifier
	counter  *countHolder
}

// NewVerifier creates a new Verifier with the associated VerifySettings settings
func NewVerifier(t testingT, settings VerifySettings) Verifier {
	if st, ok := settings.(*verifySettings); ok {
		return &verifier{
			settings: st,
			inner:    newInnerVerifier(t, st),
		}
	}

	panic("Use `NewSettings` function to create the settings.")
}

func isNil(target interface{}) bool {
	return target == nil ||
		(reflect.ValueOf(target).Kind() == reflect.Ptr && reflect.ValueOf(target).IsNil())
}
