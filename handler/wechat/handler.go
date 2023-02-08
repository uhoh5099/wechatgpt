package wechat

import (
	"github.com/eatmoreapple/openwechat"
	log "github.com/sirupsen/logrus"
)

type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type Type string

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

var handlers map[Type]MessageHandlerInterface

func init() {
	handlers = make(map[Type]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewUserMessageHandler()
}

func Handler(msg *openwechat.Message) {
	//err := handlers[GroupHandler].handle(msg)
	//if err != nil {
	//	log.Errorf("handle error: %s\n", err.Error())
	//	return
	//}

	// 处理群消息
	if msg.IsSendByGroup() {
		_ = handlers[GroupHandler].handle(msg)
		return
	}

	// 好友申请
	if msg.IsFriendAdd() {
		//if config.LoadConfig().AutoPass {
		_, err := msg.Agree("你好我是基于chatGPT引擎开发的微信机器人，你可以向我提问任何问题。")
		if err != nil {
			log.Fatalf("add friend agree error : %v", err)
			return
		}
		//}
	}

	// 私聊
	_ = handlers[UserHandler].handle(msg)
}
