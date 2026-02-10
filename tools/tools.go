//go:build tools
// +build tools

package tools

// We create a separate go.mod right in the tools folder in order
// to minimise the main go.mod file in the root folder to
// minimise public dependencies.
//
// The main (root) go.mod file will depend only on those packages
// required to actually run the app

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)
