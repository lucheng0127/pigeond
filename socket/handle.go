package socket

import (
	"encoding/json"
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

type resopnse struct {
	exitCode int
	stdout   string
	stderr   string
}

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
	rsp := resopnse{exitCode: 0}
	result, err := callbackFunc(received)
	if err != nil {
		log.Log.Error(err)
		rsp.exitCode = 1
		rsp.stderr = err.Error()
	} else {
		rsp.stdout = string(result)
	}

	// Send json response to conn
	rspJSONByte, err := json.Marshal(&rsp)
	if err != nil {
		log.Log.Error(serverError("unix", "JSON formate response error", err))
		conn.Close()
	}
	log.Log.Debugf("Send %s back to connection", string(rspJSONByte))
	conn.Write(rspJSONByte)
}
