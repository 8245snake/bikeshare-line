package main

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

//CommandHandler コマンドを処理
func CommandHandler(event *linebot.Event, message *linebot.TextMessage) {
	fmt.Printf("%v\n", message)
}
