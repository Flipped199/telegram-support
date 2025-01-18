package db

type TopicStatus string

const (
	Closed  TopicStatus = "Closed"
	Opened  TopicStatus = "Opened"
	Deleted TopicStatus = "Deleted"
)

type Topic struct {
	ThreadId   int   `xorm:"pk"`
	UserChatId int64 `xorm:"pk"`
	Status     TopicStatus
}

type MessageMap struct {
	GroupMsgId int   `xorm:"pk"`
	UserChatId int64 `xorm:"pk"`
	UserMsgId  int   `xorm:"pk"`
}
