package verify

// Verify.Go is a snapshot tool for Go language that simplifies the assertion of complex data models and documents.
//
// The library contains the following packages:
//
// The verifier package provides a set of APIs that tie in to the Go testing system and allow snapshot testing.
//
// The diff package contains interactions with well-known and custom diff tools.

// blank imports help docs.
import (
	// diff package
	_ "github.com/VerifyTests/Verify.Go/diff"
	// verifier package
	_ "github.com/VerifyTests/Verify.Go/verifier"
)