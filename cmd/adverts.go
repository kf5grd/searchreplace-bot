package cmd

import (
	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

func (b *bot) advertiseCommands() {
	b.log.Debug("Advertising commands")
	b.k.ClearCommands()
	opts := keybase.AdvertiseCommandsOptions{
		Advertisements: []chat1.AdvertiseCommandAPIParam{
			{
				Typ: "public",
				Commands: []chat1.UserBotCommandInput{
					{
						Name:        "sr add",
						Usage:       "<CONVERSATION>",
						Description: "Get the command to add this bot as a restricted bot to CONVERSATION, including all of its triggers",
					},
					{
						Name:        "sr edit",
						Description: "Get the command to edit this bot as a restricted bot to include all of its triggers",
					},
					{
						Name:        "sr about",
						Description: "Some information about the bot",
					},
				},
			},
		},
	}
	b.k.AdvertiseCommands(opts)
}

func (b *bot) clearCommands() {
	b.log.Debug("Clearing commands")
	b.k.ClearCommands()
}
