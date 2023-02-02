// Package verifier provides a set of APIs that tie in to the Go testing system and allow snapshot testing.
//
// # Example Usage
//
// To verify validity of a struct in a standard test file called function:
//
//	import (
//	  "testing"
//	  "github.com/VerifyTests/Verify.Go/verifier"
//	)
//
//	func TestVerifyingStructs(t *testing.T) {
//	  var person = Person{
//	    ID:         "ebced679-45d3-4653-8791-3d969c4a986c",
//	    Title:      Mr,
//	    FamilyName: "Smith",
//	    GivenNames: "John",
//	    Dob:        time.Date(2022, 02, 01, 1, 2, 3, 4, time.Local),
//	    Spouse:     "Jill",
//	    Address: Address{
//	      Street:  "4 Puddle Lane",
//	      Country: &country,
//	    },
//	    Children: []string{"Sam", "Mary"},
//	  }
//	  verifier.Verify(t, person)
//	}
//
// When the test is initially run will fail. If a diff tool is detected it will display the diff.
// To verify the result:
//   - Use the diff tool to accept the changes, or
//   - Manually copy the text to the new file
//
// After verifications, all `*.verified.*` files should be committed to source control and
// all `*.received.*` files should be excluded from source control.
//
// Verify allow configuring the verification process by using VerifyConfigure:
//
//	   import (
//	     "testing"
//	     "github.com/VerifyTests/Verify.Go/verifier"
//	   )
//
//	   func TestCustomConfiguration(t *testing.T) {
//	     verifier.NewVerifier(t,
//		   verifier.UseDirectory("../_testdata"), //use custom data directory
//		   verifier.UseExtension("ext"), //use custom extension
//		 ).Verify("Hello World!")
//	   }
//
// All assertion functions take, as the first argument, the `*testing.T` object provided by the
// testing framework. This allows the assertion funcs to write the failings and other details to
// the correct place.
//
// Every assertion function also takes an optional string message as the final argument,
// allowing custom error messages to be appended to the message the assertion method outputs.
package verifier
