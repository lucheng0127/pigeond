package socket

import (
	"io"
	"net"

	log "pigeond/log"
)

/*
Define callback type, then you can customize funcation to
deal with messages get from socket connection, then it will
return the result and send it back to socket connection.
*/
type callback func(receive []byte) ([]byte, error)

func handleUnixConn(conn *net.UnixConn, callbackFunc callback) {

	// Get data from conn
	received := make([]byte, 0)
	for {
		buf := make([]byte, 128)
		readLen, _, err := conn.ReadFromUnix(buf)
		received = append(received, buf[:readLen]...)
		if err != nil {
			if err != io.EOF {
				log.Log.Error(err)
			}
			break // EOF, all data received
		}
	}

	// Send data to callback and get result
	result, err := callbackFunc(received)
	if err != nil {
		log.Log.Error(err)
	}

	// Send result to conn
	conn.Write(result)
}
