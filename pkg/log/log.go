package log

import "github.com/sirupsen/logrus"

var logger *logrus.Logger = nil

func init() {
	logger = logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}

// SetLevel sets the log level to the given level.
//
// Supported levels are "panic", "fatal", "error", "warn",
// "warning", "info", "debug", "trace".
func SetLevel(lvl string) {
	switch lvl {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	}
}

func Println(args ...any) {
	logger.Println(args...)
}

func Printf(format string, args ...any) {
	logger.Printf(format, args...)
}

func Print(args ...any) {
	logger.Print(args...)
}

func Debugln(args ...any) {
	logger.Debugln(args...)
}

func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

func Debug(args ...any) {
	logger.Debug(args...)
}

func Infoln(args ...any) {
	logger.Infoln(args...)
}

func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

func Info(args ...any) {
	logger.Info(args...)
}

func Warnln(args ...any) {
	logger.Warnln(args...)
}

func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

func Warn(args ...any) {
	logger.Warn(args...)
}

func Errorln(args ...any) {
	logger.Errorln(args...)
}

func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

func Error(args ...any) {
	logger.Error(args...)
}

func Fatalln(args ...any) {
	logger.Fatalln(args...)
}

func Fatalf(format string, args ...any) {
	logger.Fatalf(format, args...)
}

func Fatal(args ...any) {
	logger.Fatal(args...)
}

func Panicln(args ...any) {
	logger.Panicln(args...)
}

func Panicf(format string, args ...any) {
	logger.Panicf(format, args...)
}

func Panic(args ...any) {
	logger.Panic(args...)
}

func Trace(args ...any) {
	logger.Trace(args...)
}

func Traceln(args ...any) {
	logger.Traceln(args...)
}

func Tracef(format string, args ...any) {
	logger.Tracef(format, args...)
}
