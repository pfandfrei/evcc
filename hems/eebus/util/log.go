package util

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type NopLogger struct{}

func (l *NopLogger) Printf(format string, v ...interface{}) {}

func (l *NopLogger) Println(v ...interface{}) {}
