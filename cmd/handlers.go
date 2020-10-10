package cmd

import (
	"regexp"
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
