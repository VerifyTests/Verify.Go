package verifier

import (
	"reflect"
)

type Verifier interface {
	Verify(target interface{})
}

func Verify(t testingT, target interface{}) {
	VerifyWithSetting(t, NewSettings(), target)
}

func VerifyWithSetting(t testingT, settings VerifySettings, target interface{}) {
	verifier := NewVerifier(t, settings)
	verifier.Verify(target)
}

func (v *verifier) Verify(target interface{}) {
	defer v.settings.runAfterVerify()

	if isNil(target) {
		v.assertExtensionIsNull()
		v.inner.VerifyInner("nil", nil, emptyTargets)
		return
	}

	if stringResult, ok := v.inner.TryGetToString(target); ok {
		if len(stringResult.Extension) > 0 {
			v.settings.UseExtension(stringResult.Extension)
		}

		value := stringResult.Value
		if len(value) == 0 {
			v.inner.VerifyInner("emptyString", nil, emptyTargets)
			return
		}

		v.inner.VerifyInner(value, nil, emptyTargets)
		return
	}

	if target, ok := target.([]byte); ok {
		v.inner.VerifyStream(target)
		return
	}

	v.assertExtensionIsNull()
	v.inner.VerifyInner(target, nil, emptyTargets)
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
