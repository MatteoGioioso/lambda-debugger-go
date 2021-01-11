package logger

type Port interface {
	Info(message interface{})
	Failure(err error)
}
