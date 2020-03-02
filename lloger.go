package llog

//Level print_level
type Level = uint

//loger level
const (
	DEBUG Level = iota
	INFO
	WARING
	ERROR
)

const logerMaxPrintLevl = ERROR

// Loger interface
type Loger interface {
	Debugln(...interface{})
	Infoln(...interface{})
	Waringln(...interface{})
	Errorln(...interface{})

	Debugfl(string, ...interface{})
	Infofl(string, ...interface{})
	Waringfl(string, ...interface{})
	Errorfl(string, ...interface{})
}
