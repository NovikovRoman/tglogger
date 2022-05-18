package tglogger

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Logger struct {
	bot    *bot
	prefix string
	chatID int64
	level  Level
}

type Settings struct {
	URL    string
	Token  string
	ChatID int64
	Prefix string

	// HTTP Client used to make requests to telegram api
	Client *http.Client
}

func New(pref Settings) (l *Logger, err error) {
	l = &Logger{
		chatID: pref.ChatID,
		prefix: pref.Prefix,
		level:  TraceLevel,
	}

	l.bot, err = newBot(pref)
	return
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

	if level > l.level {
		return "", nil
	}

	msg = fmt.Sprintf("%s %s", level, msg)

	if l.prefix != "" {
		msg = fmt.Sprintf("_%s:_ %s", l.prefix, msg)
	}

	if len(fields) > 0 {
		msg += fmt.Sprintf("\n```\n%s```", fields)
	}

	r := []rune(msg)
	if len(r) > 1100 {
		if err = l.bot.sendMessage(l.chatID, string(r[0:1024])+"…"); err != nil {
			return msg, err
		}

		f := &file{
			name:     time.Now().Format("2006-01-02_15:04:05") + "_full.log",
			caption:  "full log",
			reader:   strings.NewReader(msg),
			fileSize: 0,
		}
		err = l.bot.sendDocument(l.chatID, f)

	} else {
		err = l.bot.sendMessage(l.chatID, msg)
	}

	return msg, err
}
