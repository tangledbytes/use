package shutdown

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
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
			fmt.Println(err)
		}
	}
}
