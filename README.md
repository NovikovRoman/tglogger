# TgLogger

> Logging in telegram chats, channels.

## Example

```go
package main

import (
	"github.com/NovikovRoman/tglogger"
	"log"
)

func main() {
	var (
		tgLog *tglogger.Logger
		msg   string
		err   error
	)

	token := "12…:42…asd"
	chatID := int64(-10123456789)

	pref := tglogger.Settings{
		Token:  token,
		ChatID: chatID,
		Prefix: "Project name",
	}

	if tgLog, err = tglogger.New(pref); err != nil {
		log.Fatalln(err)
	}

	// You can react to an error.
	fields := tglogger.Fields{
		"userID":    4,
		"firstName": "Roman",
		"lastName":  "Novikov",
	}
	if msg, err = tgLog.Error("test error", fields); err != nil {
		log.Println(err, msg)
	}
}

```