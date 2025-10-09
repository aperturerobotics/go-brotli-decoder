//go:build deps_only
// +build deps_only

package hack

import (
	// _ imports golangci-lint
	_ "github.com/golangci/golangci-lint/v2/pkg/golinters"
	// _ imports golangci-lint commands
	_ "github.com/golangci/golangci-lint/v2/pkg/commands"
	// _ imports goimports
	_ "golang.org/x/tools/cmd/goimports"
	// _ imports gofumpt
	_ "mvdan.cc/gofumpt"
)
