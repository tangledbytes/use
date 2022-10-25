package shutdown

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/utkarsh-pro/use/pkg/log"
)

var fns []func() error

func RegisterFunc(f func() error) {
	fns = append(fns, f)
}

func OnSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	<-ch

	for _, fn := range fns {
		if err := fn(); err != nil {
			log.Errorln(err)
		}
	}
}
