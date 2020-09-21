package tasks

import (
	log "pigeond/log"
)

// TaskProxy run task with msg
func TaskProxy(msg []byte) ([]byte, error) {
	rst := make([]byte, 0)
	rst = append(rst, msg...)
	log.Log.Info(string(rst))
	return rst, nil
}
