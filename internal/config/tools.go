//go:build tools

package config

import (
	"fmt"

	_ "github.com/air-verse/air"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/magefile/mage"
	_ "github.com/princjef/gomarkdoc/cmd/gomarkdoc"
)
