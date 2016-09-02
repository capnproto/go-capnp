package rpc

import "log"

// A Logger records diagnostic information and errors that are not
// associated with a call.  The arguments passed into a log call are
// interpreted like fmt.Printf.  They should not be held onto past the
// call's return.
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type defaultLogger struct{}

func (defaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("rpc: "+format, args...)
}

func (defaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("rpc: "+format, args...)
}

func (c *Conn) infof(format string, args ...interface{}) {
	if c.log == nil {
		return
	}
	c.log.Infof(format, args...)
}

func (c *Conn) errorf(format string, args ...interface{}) {
	if c.log == nil {
		return
	}
	c.log.Errorf(format, args...)
}

// ConnLog sets the connection's log to the given Logger, which may be
// nil to disable logging.  By default, logs are sent to the standard
// log package.
func ConnLog(log Logger) ConnOption {
	return ConnOption{func(c *connParams) {
		c.log = log
	}}
}
