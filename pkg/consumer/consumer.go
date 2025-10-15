package consumer

import (
	"fmt"
	httpserver "sshmonitor/pkg/http"
)

func Consume() {
	for {
		select {
		case x := <-httpserver.Catcher.CommandChan:
			fmt.Println("message is", x)
		default:
		}
	}
}
