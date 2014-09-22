package openzwave

// A logger interface. Modelled on github.com/juju/loggo so that can be used substituted by default.
type Logger interface {
     	// Log an info message.
	Infof(message string, args ...interface{})
	// Log a warning message.
	Warningf(message string, args ...interface{})
	// Log an error message.
	Errorf(message string, args ...interface{})
	// Log a debug message.
	Debugf(message string, args ...interface{})
	// Log a trace message.
	Tracef(message string, args ...interface{})
}
