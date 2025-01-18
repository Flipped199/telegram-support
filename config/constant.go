package config

const ContextTopicKey = "Topic"
const ContextReceiverKey = "Receiver"

type Receiver string

const (
	ReceiverUser  Receiver = "User"
	ReceiverAdmin Receiver = "Admin"
)
