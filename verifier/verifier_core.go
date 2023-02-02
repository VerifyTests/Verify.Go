package verifier

import (
	"io"
	"reflect"
)

// Verifier is the main interface for verification process.
type Verifier interface {
	Verify(target interface{})
	Configure(configure ...VerifyConfigure) Verifier
}

// Verify verifies the passed target with the default settings.
func Verify(t testingT, target interface{}) {
	v := NewVerifier(t, nil)
	v.Verify(target)
}

// Verify verifies the passed target with the associated settings
func (v *verifier) Verify(target interface{}) {

	inner := createInnerVerifier(v.settings.t, v.settings)

	defer v.settings.runAfterVerify()

	if isNil(target) {
		v.assertExtensionIsNull()
		inner.verifyInner("nil", nil, emptyTargets)
		return
	}

	if stringResult, ok := inner.tryGetToString(target); ok {
		if len(stringResult.Extension) > 0 {
			v.settings.extension = stringResult.Extension
		}

		value := stringResult.Value
		if len(value) == 0 {
			inner.verifyInner("emptyString", nil, emptyTargets)
			return
		}

		inner.verifyInner(value, nil, emptyTargets)
		return
	}

	if target, ok := target.(io.Reader); ok {
		inner.verifyReader(target)
		return
	}

	if target, ok := target.([]byte); ok {
		inner.verifyStream(target, "bin")
		return
	}

	v.assertExtensionIsNull()
	inner.verifyInner(target, nil, emptyTargets)
}

func (v *verifier) assertExtensionIsNull() {
	if v.settings.extension == "" {
		return
	}

	panic("`UseExtension` should only be used for text, for streams, or for converter discovery.\nWhen serializing an instance the default is txt.\nTo use json as an extension when serializing use `UseStrictJSON`")
}

type verifier struct {
	settings *verifySettings
	counter  *countHolder
}

// Configure further configures the verifier
func (v *verifier) Configure(configure ...VerifyConfigure) Verifier {
	for _, cfg := range configure {
		cfg(v.settings)
	}
	return v
}

// NewVerifier creates a new Verifier with the associated settings
func NewVerifier(t testingT, configure ...VerifyConfigure) Verifier {

	var settings = newSettings(t)

	for _, cfg := range configure {
		cfg(settings)
	}

	return &verifier{
		settings: settings,
	}
}

func isNil(target interface{}) bool {
	return target == nil ||
		(reflect.ValueOf(target).Kind() == reflect.Ptr && reflect.ValueOf(target).IsNil())
}
