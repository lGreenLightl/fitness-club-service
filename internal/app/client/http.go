package client

import (
	"net"
	"time"
)

func isPortAvailable(timeout time.Duration, address string) bool {
	availableChan := make(chan struct{})
	timeoutChan := time.After(timeout)

	go func() {
		for {
			select {
			case <-timeoutChan:
				return
			default:
			}

			_, err := net.Dial("tcp", address)
			if err == nil {
				close(availableChan)
				return
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {
	case <-availableChan:
		return true
	case <-timeoutChan:
		return false
	}
}
