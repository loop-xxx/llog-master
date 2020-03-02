package lloger

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/llog-master/lloger/lprinter"
	"github.com/llog-master/lloger/lutils"
)

const minSize = int64(0x1000)
const minCheckduration = time.Second

type loginfo struct {
	format   string
	prefix   string
	logtime  time.Time
	objs     []interface{}
	isFormat bool
}

// CommonLoger ...
type CommonLoger struct {
	plevel        Level
	splitsize     int64
	logdir        string
	logname       string
	lasttime      time.Time
	checkduration time.Duration
	writer        lprinter.Printer

	lipchan chan *loginfo
	wg      sync.WaitGroup
}

func (lp *CommonLoger) installprinter() {
	if core, err := os.OpenFile(path.Join(lp.logdir, lutils.GenerateRunLogFileName(lp.logname)), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err == nil {
		lp.writer = lprinter.Build(core)
	} else {
		panic(err)
	}
}

func (lp *CommonLoger) replaceprinter() {
	if core, err := os.OpenFile(path.Join(lp.logdir, lutils.GenerateRunLogFileName(lp.logname)), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err == nil {
		lp.writer = lprinter.Build(core)
	} else {
		panic(err)
	}
}

func (lp *CommonLoger) movelog(logtime time.Time) {
	if err := os.Rename(path.Join(lp.logdir, lutils.GenerateRunLogFileName(lp.logname)), lutils.GenerateLogFileName(lp.logdir, lp.logname, lp.lasttime, logtime)); err == nil {
		lp.lasttime = logtime
	} else {
		panic(err)
	}
}

func (lp *CommonLoger) splitlog(logtime time.Time) {
	if logtime.Sub(lp.lasttime) >= lp.checkduration {
		if fi, _ := lprinter.Info(lp.writer); fi.Size() >= lp.splitsize {
			//close old writer
			lprinter.Close(lp.writer)
			//split log
			lp.movelog(logtime)
			//replace new writer
			lp.replaceprinter()
		}
	}

	return
}

func (lp *CommonLoger) end() {
	//close current wirter
	lprinter.Close(lp.writer)
	//split log
	lp.movelog(time.Now())
}

func (lp *CommonLoger) println(logtime time.Time, prefix string,
	obj ...interface{}) {
	lp.splitlog(logtime)
	lp.writer.Println(prefix, obj...)
}

func (lp *CommonLoger) printfl(logtime time.Time, prefix string,
	format string, obj ...interface{}) {
	lp.splitlog(logtime)
	lp.writer.Printfl(prefix, format, obj...)
}

//logInfo pointer channel is closed, and the data in the logInfo pointer channel is processed, worker exits
func (lp *CommonLoger) worker() {
	defer lp.wg.Done()
	for lip := range lp.lipchan {
		if lip.isFormat {
			lp.printfl(lip.logtime, lip.prefix, lip.format, lip.objs...)
		} else {
			lp.println(lip.logtime, lip.prefix, lip.objs...)
		}
	}
	lp.end()
}

/**
*	plevel: 限制打印等级, 只有>=该等级的log才会被打印
*	logdir: 日志文件存放的目录名
*	logname: 日志名
*	cpacity: loginfo buf capacity
*	splitsize: 设置log日志文件被分割的最小大小, 默认为10mb
*	checkduration: 设置每次检测日志文件大小的间隔, 默认为1s
*	日志文件最终被分割为多大由 splitsize 和 checkduration 和日志的打印速度决定
 */

//Build ...
func Build(plevel Level, logdir string, logname string, capacity uint,
	splitsize int64, checkduration time.Duration) (lp *CommonLoger, err error) {
	//plevel default = DEBUG
	if plevel > logerMaxPrintLevl {
		plevel = DEBUG
		fmt.Println("[CommonLoger WARING]:plevel out of range, set plevel = DEBUG")
	}
	if splitsize < minSize {
		splitsize = minSize
		fmt.Printf("[CommonLoger WARING]:splitsize < %v, set splitsize = %v\n", minSize, minSize)
	}
	if checkduration < minCheckduration {
		checkduration = minCheckduration
		fmt.Printf("[CommonLoger WARING]:check duration < %v, set check duration = %v\n", minCheckduration, minCheckduration)
	}

	//splitlog logdir
	err = errors.New("logdir is invaild")
	if lutils.DirIsValid(logdir) {

		//splitlog logname
		err = errors.New("log name is invaild")
		if lutils.FileNameIsValid(logname) {

			err = nil
			lp = &CommonLoger{
				plevel:        plevel,
				logdir:        logdir,
				logname:       logname,
				splitsize:     splitsize,
				checkduration: checkduration,
				lasttime:      time.Now(),
				lipchan:       make(chan *loginfo, capacity),
			}
			//install pointer
			lp.installprinter()

			//run worker
			lp.wg.Add(1)
			go lp.worker()
		}
	}

	return
}

//Close ...
func Close(lp *CommonLoger) {
	//close loginfo pointer channel
	close(lp.lipchan)
	//wait worker exit
	lp.wg.Wait()

	fmt.Println("[ConmonLoger INFO]:closed...")
}

//Debugln ...
func (lp *CommonLoger) Debugln(objs ...interface{}) {
	if DEBUG >= lp.plevel {
		logtime := time.Now()
		select {
		case lp.lipchan <- &loginfo{
			isFormat: false,
			prefix:   lutils.GenerateLogPrefix("DEBUG", logtime),
			objs:     objs,
			logtime:  logtime,
		}:
		default:
			//Don't block the logic code. If the LogInfo pointer Chan is full, discard the loginfo directly
		}
	}
}

//Infoln ...
func (lp *CommonLoger) Infoln(objs ...interface{}) {
	if INFO >= lp.plevel {
		logtime := time.Now()
		select {
		case lp.lipchan <- &loginfo{
			isFormat: false,
			prefix:   lutils.GenerateLogPrefix("INFO", logtime),
			objs:     objs,
			logtime:  logtime,
		}:
		default:
		}
	}

}

//Waringln ...
func (lp *CommonLoger) Waringln(objs ...interface{}) {
	if WARING >= lp.plevel {

		logtime := time.Now()
		select {
		case lp.lipchan <- &loginfo{
			isFormat: false,
			prefix:   lutils.GenerateLogPrefix("WARING", logtime),
			objs:     objs,
			logtime:  logtime,
		}:
		default:
		}
	}
}

//Errorln ...
func (lp *CommonLoger) Errorln(objs ...interface{}) {
	if ERROR >= lp.plevel {

		logtime := time.Now()
		lp.lipchan <- &loginfo{
			isFormat: false,
			prefix:   lutils.GenerateLogPrefix("ERROR", logtime),
			objs:     objs,
			logtime:  logtime,
		}
	}
}

//Debugfl ...
func (lp *CommonLoger) Debugfl(format string, objs ...interface{}) {
	if DEBUG >= lp.plevel {
		logtime := time.Now()
		select {
		case lp.lipchan <- &loginfo{
			isFormat: true,
			format:   format,
			prefix:   lutils.GenerateLogPrefix("DEBUG", logtime),
			objs:     objs,
			logtime:  logtime,
		}:
		default:
		}

	}

}

//Infofl ...
func (lp *CommonLoger) Infofl(format string, objs ...interface{}) {
	if INFO >= lp.plevel {

		logtime := time.Now()
		select {
		case lp.lipchan <- &loginfo{
			isFormat: true,
			format:   format,
			prefix:   lutils.GenerateLogPrefix("INFO", logtime),
			objs:     objs,
			logtime:  logtime,
		}:
		default:
		}
	}
}

//Waringfl ...
func (lp *CommonLoger) Waringfl(format string, objs ...interface{}) {
	if WARING >= lp.plevel {

		logtime := time.Now()
		select {
		case lp.lipchan <- &loginfo{
			isFormat: true,
			format:   format,
			prefix:   lutils.GenerateLogPrefix("WARING", logtime),
			objs:     objs,
			logtime:  logtime,
		}:
		default:
		}
	}
}

//Errorfl ...
func (lp *CommonLoger) Errorfl(format string, objs ...interface{}) {
	if ERROR >= lp.plevel {

		logtime := time.Now()
		lp.lipchan <- &loginfo{
			isFormat: true,
			format:   format,
			prefix:   lutils.GenerateLogPrefix("ERROR", logtime),
			objs:     objs,
			logtime:  logtime,
		}
	}
}
