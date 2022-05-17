package tglogger

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/telebot.v3"
)

type Logger struct {
	teleBot   *telebot.Bot
	botAPI    *tgbotapi.BotAPI
	prefix    string
	channelID int64
	level     Level
}

// NewTeleBotLogger https://github.com/tucnak/telebot
func NewTeleBotLogger(bot *telebot.Bot, channelID int64) *Logger {
	return &Logger{
		teleBot:   bot,
		channelID: channelID,
	}
}

// NewBotApiLogger https://github.com/go-telegram-bot-api/telegram-bot-api
func NewBotApiLogger(bot *tgbotapi.BotAPI, channelID int64) *Logger {
	return &Logger{
		botAPI:    bot,
		channelID: channelID,
	}
}

// SetPrefix set message prefix
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

// GetPrefix returns the message prefix
func (l *Logger) GetPrefix() string {
	return l.prefix
}

// SetLevel set logger level
func (l *Logger) SetLevel(lvl Level) {
	l.level = lvl
}

// GetLevel get logger level
func (l *Logger) GetLevel() Level {
	return l.level
}

func (l *Logger) Panic(msg string, fields Fields) (string, error) {
	return l.send(PanicLevel, msg, fields)
}

func (l *Logger) Fatal(msg string, fields Fields) (string, error) {
	return l.send(FatalLevel, msg, fields)
}

func (l *Logger) Error(msg string, fields Fields) (string, error) {
	return l.send(ErrorLevel, msg, fields)
}

func (l *Logger) Warn(msg string, fields Fields) (string, error) {
	return l.send(WarnLevel, msg, fields)
}

func (l *Logger) Info(msg string, fields Fields) (string, error) {
	return l.send(InfoLevel, msg, fields)
}

func (l *Logger) Log(msg string, fields Fields) (string, error) {
	return l.send(InfoLevel, msg, fields)
}

func (l *Logger) Debug(msg string, fields Fields) (string, error) {
	return l.send(DebugLevel, msg, fields)
}

func (l *Logger) Trace(msg string, fields Fields) (string, error) {
	return l.send(TraceLevel, msg, fields)
}

func (l *Logger) send(level Level, msg string, fields Fields) (string, error) {
	var err error

	if level < l.level {
		return "", nil
	}

	msg = level.String() + " " + msg

	if l.prefix != "" {
		msg = l.prefix + ": " + msg
	}

	if len(fields) > 0 {
		msg += fmt.Sprintf("```\n%s```", fields)
	}

	if l.botAPI != nil {
		m := tgbotapi.NewMessage(l.channelID, msg)
		m.ParseMode = tgbotapi.ModeMarkdown
		_, err = l.botAPI.Send(m)

	} else if l.teleBot != nil {
		chat := &telebot.Chat{
			ID:   l.channelID,
			Type: telebot.ChatChannel,
		}
		_, err = l.teleBot.Send(chat, msg, telebot.ModeMarkdown)
	}

	return msg, err
}
