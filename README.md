# TgLogger

> Логирование в telegram канал.

```go
package main

import (
	"github.com/NovikovRoman/tglogger"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/telebot.v3"
	"log"
	"time"
)

func main() {
	var (
		teleBot *telebot.Bot
		botAPI  *tgbotapi.BotAPI
		msg     string
		err     error
	)

	token := "12…:42…asd"
	channelID := int64(-10123456789)

	// Отправка с помощью gopkg.in/telebot.v3
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	if teleBot, err = telebot.NewBot(pref); err != nil {
		log.Fatalln(err)
	}

	tgLog := tglogger.NewTeleBotLogger(teleBot, channelID)

	// Можно отреагировать на ошибку отправки.
	if msg, err = tgLog.Error("test error teleBot", nil); err != nil {
		log.Println(err, msg)
	}

	// Отправка с помощью github.com/go-telegram-bot-api/telegram-bot-api/v5
	if botAPI, err = tgbotapi.NewBotAPI(token); err != nil {
		log.Fatalln(err)
	}

	tgLog = tglogger.NewBotApiLogger(botAPI, channelID)
	// Можно отреагировать на ошибку отправки.
	if msg, err = tgLog.Error("test error botAPI", nil); err != nil {
		log.Println(err, msg)
	}
}

```