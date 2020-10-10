package cmd

import (
	"searchreplacebot/pkg/logr"

	"samhofi.us/x/keybase/v2"
)

type bot struct {
	k              *keybase.Keybase
	handlers       keybase.Handlers
	opts           keybase.RunOptions
	log            *logr.Logger
	replacersBasic []string
	replacersRegex []string
}
