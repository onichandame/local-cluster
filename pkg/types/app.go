package types

import "github.com/hashicorp/go-version"

type AppDefinition struct {
	Version version.Version
	Name    string
	Command string
	Entry   string
}
