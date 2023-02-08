package wechat

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/wechatgpt/wechatbot/config"
	"github.com/wechatgpt/wechatbot/openai"
	"github.com/wechatgpt/wechatbot/utils"
	"log"
	"strings"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收私聊消息
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)

	//// 向GPT发起请求
	//requestText := strings.TrimSpace(msg.Content)
	//requestText = strings.Trim(msg.Content, "\n")
	//reply, err := gtp.Completions(requestText)
	//if err != nil {
	//	log.Printf("gtp request error: %v \n", err)
	//	msg.ReplyText("机器人神了，我一会发现了就去修。")
	//	return err
	//}
	//if reply == "" {
	//	return nil
	//}
	//
	//// 回复用户
	//reply = strings.TrimSpace(reply)
	//reply = strings.Trim(reply, "\n")
	//_, err = msg.ReplyText(reply)
	//if err != nil {
	//	log.Printf("response user error: %v \n", err)
	//}
	if msg.IsSendBySelf() {
		msg.FromUserName = msg.ToUserName
	}

	wechat := config.GetWechatKeyword()
	requestText := msg.Content
	if wechat != nil {
		content, key := utils.ContainsI(requestText, *wechat)
		if len(key) == 0 {
			return nil
		}
		splitItems := strings.Split(content, key)
		if len(splitItems) < 2 {
			return nil
		}
		requestText = strings.TrimSpace(splitItems[1])
	}
	log.Println("问题：", requestText)
	reply, err := openai.Completions(requestText)
	if err != nil {
		log.Println(err)
		if reply != nil {
			result := *reply
			// 如果文字超过4000个字会回错，截取前4000个文字进行回复
			if len(result) > 4000 {
				_, err = msg.ReplyText(result[:4000])
				if err != nil {
					log.Println("回复出错：", err.Error())
					return err
				}
			}
		}
		text, err := msg.ReplyText(fmt.Sprintf("bot error: %s", err.Error()))
		log.Println(text)
		return err
	}
	// 如果在提问的时候没有包含？,AI会自动在开头补充个？看起来很奇怪
	result := *reply
	if strings.HasPrefix(result, "?") {
		result = strings.Replace(result, "?", "", -1)
	}
	if strings.HasPrefix(result, "？") {
		result = strings.Replace(result, "？", "", -1)
	}
	// 微信不支持markdown格式，所以把反引号直接去掉
	if strings.Contains(result, "`") {
		result = strings.Replace(result, "`", "", -1)
	}

	if reply != nil {
		_, err = msg.ReplyText(*reply)
		if err != nil {
			log.Println(err)
		}
		return err
	}
	return err
}
