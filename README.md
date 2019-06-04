一个酷Q-HTTP API插件的，支持多客户端的，反向websocket服务[开发中]
========


例子：

```
// 消息处理函数
func handleMessage(client *cqws.Client, message *models.Message) {
	fmt.Printf("%+v\n", message)
	fmt.Println("机器人QQ:", client.QQ)
	// 私聊消息回复
	if message.MessageType == "private" {
		_ = client.SendMessage("send_private_msg", map[string]interface{}{
			"user_id": message.UserID,
			"message": "收到",
		})
		return
	}
	// 群消息回复
	if message.MessageType == "group" {
		_ = client.SendMessage("send_group_msg", map[string]interface{}{
			"group_id":    message.UserID,
			"message":     "收到",
			"auto_escape": true,
		})
		return
	}
}

func main() {
	// 更多配置参数见sopt模块下面的函数，也可以自己定义opt函数，NewServer会依次调用opt函数
	server := cqws.NewServer(
		sopt.MessageHandler(handleMessage),
	)
	server.Listen()
}
```