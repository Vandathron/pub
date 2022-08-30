package pub

type Logger interface {
	LogInfo(msg string)
	LogErr(msg string)
}

