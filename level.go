package tglogger

// Level type
type Level uint32

const (
	PanicLevel = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	switch level {
	case 0:
		return "🆘"
	case 1:
		return "❌"
	case 2:
		return "❗"
	case 3:
		return "⚠"
	case 4:
		return "🗒"
	case 5:
		return "📝"
	case 6:
		return "📜"
	}

	return "⁉"
}
