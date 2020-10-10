package cmd

import (
	"fmt"
	"regexp"
	"searchreplacebot/pkg/util"
	"strings"

	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
	"samhofi.us/x/keybase/v2/types/stellar1"
)

func (b *bot) registerHandlers() {
	b.log.Debug("Registering handlers")

	var (
		chat   = b.chatHandler
		conv   = b.convHandler
		wallet = b.walletHandler
		err    = b.errorHandler
	)
	b.handlers = keybase.Handlers{
		ChatHandler:         &chat,
		ConversationHandler: &conv,
		WalletHandler:       &wallet,
		ErrorHandler:        &err,
	}
}

func (b *bot) chatHandler(m chat1.MsgSummary) {
	var (
		userName = m.Sender.Username
		convID   = m.ConvID
		msgID    = m.Id
	)

	if userName == b.k.Username {
		return
	}

	if len(b.filterConvs) > 0 && !util.ConvIDInSlice(convID, b.filterConvs) {
		return
	}

	switch m.Content.TypeName {
	case "text":
		text := m.Content.Text.Body

		if strings.HasPrefix(text, "!sr add") || strings.HasPrefix(text, "!sr edit") {
			var triggers = make([]string, 0)
			for _, r := range b.replacersRegex {
				if len(r) < 4 {
					continue
				}
				separator := string([]rune(r)[0])
				replacer := string([]rune(r)[1:])
				s := strings.Split(replacer, separator)
				if len(s) < 2 {
					continue
				}
				_, err := regexp.Compile(s[0])
				if err != nil {
					continue
				}
				triggers = append(triggers, "--allow-trigger '"+strings.Replace(s[0], "'", `\'`, -1)+"'")
			}
			for _, r := range b.replacersBasic {
				if len(r) < 4 {
					continue
				}
				separator := string([]rune(r)[0])
				replacer := string([]rune(r)[1:])
				s := strings.Split(replacer, separator)
				if len(s) < 2 {
					continue
				}
				triggers = append(triggers, "--allow-trigger '"+strings.Replace(s[0], "'", `\'`, -1)+"'")
			}
			if strings.HasPrefix(text, "!sr add") {
				msg := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!sr add", "", 1))
				if msg == "" {
					return
				}
				conversation := strings.Fields(msg)[0]
				b.k.SendMessageByConvID(convID, "```keybase chat add-bot-member -u '"+b.k.Username+"' -r 'restrictedbot' --allow-commands "+strings.Join(triggers, " ")+" '"+conversation+"'```")
				return
			}
			if strings.HasPrefix(text, "!sr edit") {
				b.k.SendMessageByConvID(convID, "```keybase chat edit-bot-member -u '"+b.k.Username+"' -r 'restrictedbot' --allow-commands "+strings.Join(triggers, " ")+" '"+m.Channel.Name+"'```")
				return
			}
		}

		if strings.HasPrefix(text, "!sr about") {
			var donations string
			backTick := "`"

			if b.k.Username == "dont_furl_me_bro" {
				donations = `

Server resources are very limited around here, and any donations you'd like to send me are greatly appreciated! If you'd like to make a donation to help offset server costs, you can do so by sending me XLM from your Keybase wallet. Any amount helps!
`
			}
			aboutMsg := fmt.Sprintf(`*SearchReplaceBot* _ by @dxb _

This is a fairly simple bot which watches for certain text, replaces it with something else, and replies with the result. If you'd like to run your own instance of the bot you can download its source code at https://github.com/kf5grd/searchreplace-bot.

If you'd like to add this particular instance of the bot to your own conversation, send me the command %s!sr add <CONVERSATION>%s, where %s<CONVERSATION>%s is the name of the team, or PM you'd like to add me to, then copy the command I send you and paste it into your terminal.%s
`, backTick, backTick, backTick, backTick, donations)

			b.k.SendMessageByConvID(convID, aboutMsg)
			return
		}

		for _, r := range b.replacersRegex {
			if len(r) < 4 {
				continue
			}
			separator := string([]rune(r)[0])
			replacer := string([]rune(r)[1:])
			s := strings.Split(replacer, separator)
			if len(s) < 2 {
				continue
			}
			regx, err := regexp.Compile(s[0])
			if err != nil {
				continue
			}
			text = regx.ReplaceAllString(text, s[1])
		}
		if text != m.Content.Text.Body {
			text = strings.Replace(text, "%", "%%", -1)
			b.k.ReplyByConvID(convID, msgID, text)
			return
		}

		for _, r := range b.replacersBasic {
			if len(r) < 4 {
				continue
			}
			separator := string([]rune(r)[0])
			replacer := string([]rune(r)[1:])
			s := strings.Split(replacer, separator)
			if len(s) < 2 {
				continue
			}
			text = strings.Replace(text, s[0], s[1], -1)
		}
		if text == m.Content.Text.Body {
			return
		}
		text = strings.Replace(text, "%", "%%", -1)
		b.k.ReplyByConvID(convID, msgID, text)
		return
	}
}

func (b *bot) convHandler(c chat1.ConvSummary) {
}

func (b *bot) walletHandler(p stellar1.PaymentDetailsLocal) {
}

func (b *bot) errorHandler(e error) {
	b.log.Error("Error: %v", e)
}
