package middleware

import (
	"github.com/Flipped199/telegram-support/config"
	"github.com/Flipped199/telegram-support/db"
	"github.com/Flipped199/telegram-support/handler"
	tele "gopkg.in/telebot.v4"
)

func CheckUserStatus(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		cfg := config.GetConfig()
		// 如果是用户侧消息，检查用户状态
		var topic db.Topic
		var receiver config.Receiver
		if c.Chat().Type == tele.ChatPrivate {
			receiver = config.ReceiverAdmin
			//if !slices.Contains(cfg.Admin, c.Chat().ID) {
			topic.UserChatId = c.Chat().ID
			get, err := db.GetDB().Get(&topic)
			if err != nil {
				return err
			}
			if !get {
				// 创建一个topic
				newTopic, err := handler.CreateTopicAndSave(c)
				if err != nil {
					return err
				}
				topic.ThreadId = newTopic.ThreadID
				topic.Status = db.Opened
			} else {
				switch topic.Status {
				case db.Closed:
					// 重新打开topic
					if err = c.Bot().ReopenTopic(&tele.Chat{ID: cfg.GroupId},
						&tele.Topic{ThreadID: topic.ThreadId}); err != nil {
						return err
					}
					topic.Status = db.Opened
					if _, err = db.GetDB().Update(topic); err != nil {
						return err
					}
				case db.Deleted:
					// 提示用户
					return c.Send("您好，您似乎被禁止联络管理员，请尝试其它途径。")
				}
			}
		} else if c.Chat().ID == cfg.GroupId && c.Message().ThreadID != 0 { // 如果是管理在群内发送的消息
			receiver = config.ReceiverUser

			topic.ThreadId = c.Message().ThreadID
			get, err := db.GetDB().Get(&topic)
			if err != nil {
				return err
			}
			if !get {
				return c.Send("当前话题没有绑定的用户，这可能是由于您手动创建了该话题")
			}
		} else {
			return nil
		}
		c.Set(config.ContextTopicKey, &topic)
		c.Set(config.ContextReceiverKey, receiver)

		return next(c)
	}
}
