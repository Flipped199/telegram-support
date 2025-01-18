package handler

import (
	"fmt"
	"github.com/Flipped199/telegram-support/config"
	"github.com/Flipped199/telegram-support/db"
	"github.com/Flipped199/telegram-support/log"
	tele "gopkg.in/telebot.v4"
	"log/slog"
)

func OnStart(c tele.Context) error {
	// 查话题映射
	topic := c.Get(config.ContextTopicKey).(*db.Topic)

	log.Debug("start bot", slog.Any("topic", topic))

	return c.Send("您好，请问有什么可以帮助您的？向我发送的每一条消息都将被人工处理！")
}

func OnDelCmd(c tele.Context) error {
	log.Debug("del msg")
	msg := c.Message()
	if msg.ReplyTo.ID == msg.ReplyTo.ThreadID {
		return c.Reply("tips: 请指定需要删除的消息（清空使用/clear）")
	}
	topic := c.Get(config.ContextTopicKey).(*db.Topic)
	msgMap := db.MessageMap{
		GroupMsgId: msg.ReplyTo.ID,
		UserChatId: topic.UserChatId,
	}
	get, err := db.GetDB().Get(&msgMap)
	if err != nil {
		return err
	}
	if !get {

	}
	if err = c.Bot().Delete(&tele.Message{
		ID: msgMap.UserMsgId,
		Chat: &tele.Chat{
			ID: topic.UserChatId,
		},
	}); err != nil {
		return err
	}
	return c.Bot().Delete(msg.ReplyTo)
}

func createTopic(c tele.Context) (*tele.Topic, error) {
	topic, err := c.Bot().CreateTopic(&tele.Chat{ID: config.GetConfig().GroupId}, &tele.Topic{
		Name: fmt.Sprintf("%s|%d", c.Sender().FirstName, c.Sender().ID),
	})
	if err != nil {
		return nil, fmt.Errorf("create topic: %w", err)
	}
	return topic, nil
}

func CreateTopicAndSave(c tele.Context) (*tele.Topic, error) {
	topic, err := createTopic(c)
	if err != nil {
		return nil, fmt.Errorf("create topic: %w", err)
	}
	threads := db.Topic{
		UserChatId: c.Chat().ID,
		ThreadId:   topic.ThreadID,
		Status:     db.Opened,
	}
	_, err = db.GetDB().Insert(threads)
	if err != nil {
		return nil, fmt.Errorf("insert topic: %w", err)
	}
	return topic, nil
}
