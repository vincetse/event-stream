package amqp

import (
	"time"
	log "github.com/sirupsen/logrus"
	mq "github.com/streadway/amqp"
)

func Dial(uri string) *mq.Connection {
	var conn *mq.Connection
	for i := 0; conn == nil; i++ {
		_conn, err := mq.Dial(uri)
		if err == nil {
			log.Infof("connected to %s", uri)
			conn = _conn
		} else {
			var sleepTime time.Duration = time.Duration(i * i) * 1e9
			log.Info(err)
			log.Infof("retrying in %d nanoseconds", sleepTime)
			time.Sleep(sleepTime)
		}
	}
	return conn
}
