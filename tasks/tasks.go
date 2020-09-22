package tasks

import (
	"strconv"
	"strings"

	log "pigeond/log"
)

// TaskProxy run task with msg
func TaskProxy(msg []byte) ([]byte, error) {

	// msg: [Auto ACK] [Command] [Arg1] [Arg2] ... [END]
	// eg: true ADD_SCRIPT test+script /tmp/test.tar END
	// need replace "+" to space in args
	rst := make([]byte, 0)
	rstChan := make(chan string)
	errChan := make(chan string)
	msgList := strings.Split(string(msg), " ")
	log.Log.Debug("Command:", msgList)

	if len(msgList) < 3 {
		return rst, wrongCommandError()
	}
	autoAck, err := strconv.ParseBool(msgList[0])
	if err != nil {
		return rst, err
	}
	command := msgList[1]
	args := []string{}
	if len(msgList) > 3 {
		args = msgList[2 : len(msgList)-2]
	}

	// Goroutine run task, if auto ack do not wait task finish
	switch command {
	case "LIST_SCRIPT":
		go listScripts(rstChan, errChan)
	case "ADD_SCRIPT":
		log.Log.Info("Try to add script with args:", args)
	default:
		return rst, unsupportCommandError()
	}

	// Auto Ack, don't check task result
	if autoAck == true {
		return []byte("Task Auto Ack"), nil
	}

	// Block until get result or error from channel
	select {
	case errStr := <-errChan:
		return rst, taskRunError(errStr)
	case rstStr := <-rstChan:
		rst = append(rst, []byte(rstStr)...)
		return rst, nil
	}
}
