//go:build tools

// Package tools pins runtime dependencies that are only imported by gitignored
// generated files. Without this file, Dependabot removes them because it
// cannot see the generated source files.
package tools

import (
	// Imported by oapi-codegen generated server/client code (gitignored).
	_ "github.com/apapsch/go-jsonmerge/v2"
	_ "github.com/getkin/kin-openapi/openapi3"
	_ "github.com/go-chi/chi/v5"
)
