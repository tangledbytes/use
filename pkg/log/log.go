package log

import "github.com/sirupsen/logrus"

func Println(args ...any) {
	logrus.Println(args...)
}

func Printf(format string, args ...any) {
	logrus.Printf(format, args...)
}

func Print(args ...any) {
	logrus.Print(args...)
}

func Debugln(args ...any) {
	logrus.Debugln(args...)
}

func Debugf(format string, args ...any) {
	logrus.Debugf(format, args...)
}

func Debug(args ...any) {
	logrus.Debug(args...)
}

func Infoln(args ...any) {
	logrus.Infoln(args...)
}

func Infof(format string, args ...any) {
	logrus.Infof(format, args...)
}

func Info(args ...any) {
	logrus.Info(args...)
}

func Warnln(args ...any) {
	logrus.Warnln(args...)
}
