package types

import "github.com/hashicorp/go-version"

type platformSpecificDef struct {
	Platform     Platform
	Url          string
	CompressType CompressType
}

type AppDefinition struct {
	Version version.Version
	Name    string
	Specs   []platformSpecificDef
}
