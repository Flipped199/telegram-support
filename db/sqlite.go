package db

import (
	"github.com/Flipped199/telegram-support/log"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func init() {
	var err error
	engine, err = xorm.NewEngine("sqlite3", "./bot.db")
	if err != nil {
		log.Fatal("sqlite3", err)
	}
	err = engine.Sync(new(MessageMap), new(Topic))
	if err != nil {
		log.Fatal("sqlite3 sync", err)
	}
}

func GetDB() *xorm.Engine {
	return engine
}
