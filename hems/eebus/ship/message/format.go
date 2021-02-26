package message

import "github.com/thoas/go-funk"

const (
	ProtocolHandshakeFormatJSON = "JSON-UTF8"
)

// IsSupported validates if format is supported
func (m MessageProtocolFormatsType) IsSupported(format string) bool {
	return funk.ContainsString(m.Format, format)
}
