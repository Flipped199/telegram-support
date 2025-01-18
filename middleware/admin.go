package middleware

import (
	"github.com/Flipped199/telegram-support/config"
	tele "gopkg.in/telebot.v4"
	"slices"
)

func Admin(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		cfg := config.GetConfig()
		if c.Chat().ID == cfg.GroupId && slices.Contains(cfg.Admin, c.Sender().ID) {
			return next(c)
		}

		return nil
	}
}
