package rpc

type ILogger interface {
	Debugw(mgs string, keysAndValues ...interface{})
	Infow(mgs string, keysAndValues ...interface{})
	Warnw(mgs string, keysAndValues ...interface{})
	Errorw(mgs string, keysAndValues ...interface{})
	Fatalw(mgs string, keysAndValues ...interface{})
	Panicw(mgs string, keysAndValues ...interface{})
}
