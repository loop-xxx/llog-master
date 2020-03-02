package lloger

import (
	"fmt"
	"os"
	"time"

	"github.com/llog-master/lloger/lprinter"
	"github.com/llog-master/lloger/lutils"
)

// ConsoleLoger ...
type ConsoleLoger struct {
	plevel Level
	writer lprinter.Printer
}

//BuildConsoleLoger ...
func BuildConsoleLoger(level Level) (lp *ConsoleLoger) {
	if level > logerMaxPrintLevl {
		level = DEBUG
		fmt.Println("[ConsoleLoger WARING]:level out of range, set level = DEBUG")
	}
	lp = &ConsoleLoger{plevel: level}
	lp.writer = lprinter.Build(os.Stdout)

	return
}

//Debugln ...
func (lp *ConsoleLoger) Debugln(objs ...interface{}) {
	if DEBUG >= lp.plevel {
		lp.writer.Println(lutils.GenerateLogPrefix("DEBUG", time.Now()), objs...)
	}
}

//Infoln ...
func (lp *ConsoleLoger) Infoln(objs ...interface{}) {
	if INFO >= lp.plevel {
		lp.writer.Println(lutils.GenerateLogPrefix("INFO", time.Now()), objs...)
	}
}

//Waringln ...
func (lp *ConsoleLoger) Waringln(objs ...interface{}) {
	if WARING >= lp.plevel {
		lp.writer.Println(lutils.GenerateLogPrefix("WARING", time.Now()), objs...)
	}
}

//Errorln ...
func (lp *ConsoleLoger) Errorln(objs ...interface{}) {
	if ERROR >= lp.plevel {
		lp.writer.Println(lutils.GenerateLogPrefix("ERROR", time.Now()), objs...)
	}
}

//Debugfl ...
func (lp *ConsoleLoger) Debugfl(format string, objs ...interface{}) {
	if DEBUG >= lp.plevel {
		lp.writer.Printfl(lutils.GenerateLogPrefix("DEBUG", time.Now()), format, objs...)
	}
}

//Infofl ...
func (lp *ConsoleLoger) Infofl(format string, objs ...interface{}) {
	if INFO >= lp.plevel {
		lp.writer.Printfl(lutils.GenerateLogPrefix("INFO", time.Now()), format, objs...)
	}
}

//Waringfl ...
func (lp *ConsoleLoger) Waringfl(format string, objs ...interface{}) {
	if WARING >= lp.plevel {
		lp.writer.Printfl(lutils.GenerateLogPrefix("WARING", time.Now()), format, objs...)
	}
}

//Errorfl ...
func (lp *ConsoleLoger) Errorfl(format string, objs ...interface{}) {
	if ERROR >= lp.plevel {
		lp.writer.Printfl(lutils.GenerateLogPrefix("ERROR", time.Now()), format, objs...)
	}
}
