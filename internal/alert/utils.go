package alert

const (
	TypeError   = "error"
	TypeWarning = "warning"
	TypeInfo    = "info"
	TypeSuccess = "success"
)

func addMessageTitleEmoji(messageTitle, messageType string) string {
	switch messageType {
	case TypeError:
		return "🚨 " + messageTitle
	case TypeWarning:
		return "⚠️ " + messageTitle
	case TypeInfo:
		return "ℹ️ " + messageTitle
	case TypeSuccess:
		return "✅ " + messageTitle
	default:
		return "🔔 " + messageTitle
	}
}
