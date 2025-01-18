package main

import (
	"fmt"
	"github.com/Flipped199/telegram-support/config"
	"github.com/Flipped199/telegram-support/db"
	_ "github.com/Flipped199/telegram-support/db"
	"github.com/Flipped199/telegram-support/handler"
	"github.com/Flipped199/telegram-support/log"
	"github.com/Flipped199/telegram-support/middleware"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}
	if config.GetConfig().Debug {
		log.SetLevel(slog.LevelDebug)
		log.Debug("You are running in debug mode")
	}

	startBot()
}

func startBot() {
	cfg := config.GetConfig()

	if cfg.BotToken == "" {
		log.Fatal("Bot token is empty", nil)
	}

	pref := tele.Settings{
		Token:  config.GetConfig().BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	if cfg.Proxy.URL != "" {
		proxyUrl, err := url.Parse(cfg.Proxy.URL)
		if err != nil {
			log.Fatal("Proxy URL parse error", err)
		}
		pref.Client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal("Bot initialize error", err)
	}

	initHandlers(b)
	checkPermissions(b)
	intro(b)
	b.Start()
}

func initHandlers(b *tele.Bot) {
	b.Use(middleware.CheckUserStatus)
	b.Handle("/start", handler.OnStart)
	b.Handle(tele.OnEdited, handler.OnEdited)
	b.Handle(tele.OnText, handler.OnMsg)
	b.Handle(tele.OnMedia, handler.OnMsg)
	b.Handle(tele.OnTopicClosed, func(c tele.Context) error {
		log.Debug("topic closed", slog.Any("topic", c.Message().ThreadID))
		topic := c.Get(config.ContextTopicKey).(*db.Topic)
		// 阻止用户发送消息
		topic.Status = db.Closed
		if _, err := db.GetDB().Update(topic); err != nil {
			return err
		}
		// 向用户端发送提示
		_, err := c.Bot().Send(&tele.Chat{ID: topic.UserChatId}, "会话已关闭，发送任意消息以重新接入会话。")
		return err
	})
	b.Handle(tele.OnTopicReopened, func(c tele.Context) error {
		log.Debug("topic reopened", slog.Any("topic", c.Message().ThreadID))
		topic := c.Get(config.ContextTopicKey).(*db.Topic)
		// 取消阻止用户发送消息
		topic.Status = db.Opened
		if _, err := db.GetDB().Update(topic); err != nil {
			return err
		}
		return nil
	})
	//  only admins
	admins := b.Group()
	admins.Use(middleware.Admin)
	admins.Handle("/del", handler.OnDelCmd)
}

func intro(b *tele.Bot) {
	fmt.Printf("Beep boop! [%s] is online and ready to rock!\n", b.Me.Username)
}

func checkPermissions(b *tele.Bot) {
	chat := &tele.Chat{ID: config.GetConfig().GroupId}

	member, err := b.ChatMemberOf(chat, b.Me)
	if err != nil {
		log.Fatal("Get member", err)
	}

	if !member.CanManageTopics {
		log.Fatal("The bot doesnt have permission to manage topics", nil)
	}

	sendTips(b, chat)
}

func sendTips(b *tele.Bot, chat *tele.Chat) {
	topic, err := b.CreateTopic(&tele.Chat{ID: config.GetConfig().GroupId}, &tele.Topic{
		Name: "当您看到这条消息的时候，说明您的配置完全正确",
	})
	if err != nil {
		log.Error("Create topic error", err)
	}

	_, err = b.Send(chat, "此话题将在1分钟后删除（如果你没有kill掉我的话）", &tele.SendOptions{
		ThreadID: topic.ThreadID,
	})
	if err != nil {
		log.Error("Reply message error", err)
	}

	go func() {
		<-time.After(time.Minute)
		if err := b.DeleteTopic(chat, topic); err != nil {
			log.Error("Delete topic error", err)
		}
	}()
}
