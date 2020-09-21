package socket

import (
	"fmt"
	"net"
	"os"

	log "pigeond/log"
)

// LaunchServer to start a unix socket or tcp server
func LaunchServer(serverType, socketFile, port string, callbackFunc callback) error {
	switch serverType {
	case "unix":
		// Delete unix socket file if exist
		_, err := os.Stat(socketFile)
		if err == nil {
			err = os.Remove(socketFile)
			if err != nil {
				return serverError("unix", "", err)
			}
		}

		addr, err := net.ResolveUnixAddr("unix", socketFile)
		if err != nil {
			return serverError("unix", "", err)
		}

		l, err := net.ListenUnix("unix", addr)
		if err != nil {
			return serverError("unix", "", err)
		}

		for {
			log.Log.Info("Start to listen", addr)
			conn, err := l.AcceptUnix()
			if err != nil {
				return serverError("unix", "", err)
			}

			// Call handleUnixConn to handle connection message
			go handleUnixConn(conn, callbackFunc)
		}
	case "tcp":
		fmt.Println("Start server on port", port)
		return nil
	default:
		return unsupportServerError()
	}
}
