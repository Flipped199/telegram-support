package handler

import (
	"errors"
	"github.com/Flipped199/telegram-support/config"
	"github.com/Flipped199/telegram-support/db"
	"github.com/Flipped199/telegram-support/log"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strings"
)

func OnMsg(c tele.Context) error {
	msg := c.Message()
	// 加载配置和上下文
	cfg := config.GetConfig()
	topic := c.Get(config.ContextTopicKey).(*db.Topic)
	receiver := c.Get(config.ContextReceiverKey).(config.Receiver)

	// 初始化消息发送选项
	var opts tele.SendOptions
	var targetChat tele.Recipient
	msgMap := &db.MessageMap{
		UserChatId: topic.UserChatId,
	}
	// 确定消息方向
	if receiver == config.ReceiverUser { // 客服 -> 用户
		if msg.IsReply() && msg.ReplyTo.ID != msg.ReplyTo.ThreadID {
			err := setReplyOptionsFromMessageMap(&opts, "GroupMsgId", msg.ReplyTo.ID)
			if err != nil {
				log.Warn("客服回复用户消息时未找到引用的消息")
			}
		}

		targetChat = &tele.Chat{ID: topic.UserChatId}
	} else { // 用户 -> 客服
		opts.ThreadID = topic.ThreadId

		if msg.IsReply() {
			err := setReplyOptionsFromMessageMap(&opts, "UserMsgId", msg.ReplyTo.ID)
			if err != nil {
				log.Warn("用户回复客服消息时未找到引用的消息")
			}
		}

		targetChat = &tele.Chat{ID: cfg.GroupId}
	}

	forwardMsg, err := c.Bot().Copy(targetChat, c.Message(), &opts)
	if err != nil {
		opts.ReplyTo = nil
		// 如果是回复的消息 但是引用被删除导致回复失败，则直接发送
		if strings.Contains(err.Error(), "message to be replied not found") {
			forwardMsg, err = c.Bot().Copy(targetChat, c.Message(), &opts)
			if err != nil {
				return err
			}
		} else if strings.Contains(err.Error(), "message thread not found") {
			// 话题被管理员删除，需要重新创建，话题删除后
			newTopic, err := createTopic(c)
			if err != nil {
				return err
			}
			topic.ThreadId = newTopic.ThreadID
			if _, err = db.GetDB().Update(topic); err != nil {
				return err
			}
			opts.ThreadID = newTopic.ThreadID
			// 重发消息
			forwardMsg, err = c.Bot().Copy(targetChat, c.Message(), &opts)
			if err != nil {
				return err
			}
		} else {
			log.Error("Failed to forward message", err)
			return err
		}
	}

	if receiver == config.ReceiverAdmin {
		msgMap.GroupMsgId = forwardMsg.ID
		msgMap.UserMsgId = msg.ID
		log.Info("Message forwarded to group", slog.Int("forwarded_msg_id", forwardMsg.ID))
	} else {
		msgMap.UserMsgId = forwardMsg.ID
		msgMap.GroupMsgId = msg.ID
		log.Info("Message forwarded to user", slog.Int("forwarded_msg_id", forwardMsg.ID))
	}

	_, err = db.GetDB().Insert(msgMap)
	if err != nil {
		log.Error("Failed to save message ID mapping", err)
		return err
	}

	log.Info("Message mapping saved",
		slog.Int("source_msg_id", msg.ID),
		slog.Int("mapped_msg_id", forwardMsg.ID),
	)

	return nil
}

func OnEdited(c tele.Context) error {
	log.Debug("on msg edited")
	cfg := config.GetConfig()
	msg := c.Message()
	topic := c.Get(config.ContextTopicKey).(*db.Topic)
	targetMsgId, err := getTargetMsgID(c.Message(), topic, cfg.GroupId)
	if err != nil {
		return err
	} else if targetMsgId == 0 {
		log.Debug("The message you are trying to edit doesn't exist")
		return nil
	}

	// 编辑用户端消息
	var what any
	if msg.Photo != nil {
		msg.Photo.Caption = msg.Caption
		what = msg.Photo
	} else if msg.Video != nil {
		msg.Video.Caption = msg.Caption
		what = msg.Video
	} else if msg.Document != nil {
		msg.Document.Caption = msg.Caption
		what = msg.Document
	} else if msg.Voice != nil {
		msg.Voice.Caption = msg.Caption
		what = msg.Voice
	} else if msg.Audio != nil {
		msg.Audio.Caption = msg.Caption
		what = msg.Audio
	} else {
		what = msg.Text
	}
	_, err = c.Bot().Edit(&tele.Message{
		ID: targetMsgId,
		Chat: &tele.Chat{
			ID: topic.UserChatId,
		},
	}, what)
	return err
}

// 从消息映射表中设置回复选项
func setReplyOptionsFromMessageMap(opts *tele.SendOptions, key string, msgID int) error {
	msgMap := db.MessageMap{}
	if key == "GroupMsgId" {
		msgMap.GroupMsgId = msgID
	} else if key == "UserMsgId" {
		msgMap.UserMsgId = msgID
	}

	found, err := db.GetDB().Get(&msgMap)
	if err != nil {
		return err
	}

	if found {
		if key == "GroupMsgId" {
			opts.ReplyTo = &tele.Message{ID: msgMap.UserMsgId}
		} else if key == "UserMsgId" {
			opts.ReplyTo = &tele.Message{ID: msgMap.GroupMsgId}
		}
		return nil
	}

	return errors.New("msg map not found")
}

// 从数据库中获取目标消息 ID
func getTargetMsgID(msg *tele.Message, topic *db.Topic, groupId int64) (int, error) {
	var targetMsgId int
	msgMap := db.MessageMap{
		UserChatId: topic.UserChatId,
	}

	// 根据消息来源（群组或用户）设置映射
	if msg.Chat.ID == groupId {
		msgMap.GroupMsgId = msg.ID
	} else {
		msgMap.UserMsgId = msg.ID
	}

	// 查询数据库以获取映射
	found, err := db.GetDB().Get(&msgMap)
	if err != nil {
		return 0, err
	}

	if found {
		if msg.Chat.ID == groupId {
			targetMsgId = msgMap.UserMsgId
		} else {
			targetMsgId = msgMap.GroupMsgId
		}
	}

	return targetMsgId, nil
}
