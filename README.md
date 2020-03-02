#### llog

* `llog`是一款轻量级的日志库, 支持并发,异步的记录日志。

###### 简单实用:

```go
func logerExample() {
	if ll, err := lloger.Build(lloger.DEBUG, "./varlog", "loop",
		0x80000, 10*(1<<20), 2*time.Second); err == nil {
		defer lloger.Close(ll)

		for i := uint(0); i < 0xffffff; i++ {
			ll.Debugfl("%#v", "hello, world")
			//ll.Debugln("hello, world")
		}
	} else {
		fmt.Println(err)
	}

}
```

###### 参数介绍:

```go
/**
*	plevel: 限制打印等级, 只有>=该等级的log才会被打印
*	logdir: 日志文件存放的目录名
*	logname: 日志名
*	cpacity: loginfo buf capacity
*	splitsize: 设置log日志文件被分割的最小大小, 默认为10mb
*	checkduration: 设置每次检测日志文件大小的间隔, 默认为1s
*	日志文件最终被分割为多大由 splitsize 和 checkduration 和日志的打印速度决定
 */
func Build(plevel Level, logdir string, logname string, capacity uint,
	splitsize int64, checkduration time.Duration) (lp *CommonLoger, err error)
```

###### 日志前缀格式以及生成日志文件名命名格式修改:

```go
//loop-xxx/llog-master/lutils/lutils.go
/**
日志前缀生成函数
修改日志前缀格式, 修改此函数即可
*/
const logTimeFormat = string("2006-01-02 15:04:05.000")

//GenerateLogPrefix ...
func GenerateLogPrefix(level string, nowtime time.Time) (logprefix string)
```



```go
//loop-xxx/llog-master/lutils/lutils.go
/**
日志文件名生成函数
修改命名方式, 修改此函数即可
*/
const filenameTimeformat = string("20060102150405")

//GenerateLogFileName ...
func GenerateLogFileName(dir, name string, last, now time.Time) (filename string)
/**
runlog 文件名生成函数
如果要修改runlog文件名, 修改此函数即可
*/

// GenerateRunLogFileName ...
func GenerateRunLogFileName(name string) (filename string)
```

###### console_loger支持(非异步):

```go
func consolelogerExample() {
	ll := lloger.BuildConsoleLoger(lloger.DEBUG)
	for i := uint(0); i < 0xffff; i++ {
		ll.Debugfl("%#v", "hello, world")
		ll.Debugln("hello, world")
	}
}
```

###### 并发记录log:

```go

func handlersub(ll lloger.Loger) {
	for i := uint(0); i < 0xfffff; i++ {
		ll.Debugfl("%#v", "hello, world")
		ll.Debugln("hello, world")
		ll.Infofl("%#v", "hello, world")
		ll.Infoln("hello, world")
		ll.Waringfl("%#v", "hello, world")
		ll.Waringln("hello, world")
		ll.Errorfl("%#v", "hello, world")
		ll.Errorln("hello, world")

	}
}

func handler(wgp *sync.WaitGroup, ll lloger.Loger) {
	defer wgp.Done()
	handlersub(ll)
}

func handlerpool(ll lloger.Loger) {
	var wg sync.WaitGroup

	for i := int(0); i < 0x10; i++ {
		wg.Add(1)
		go handler(&wg, ll)
	}
	handlersub(ll)

	wg.Wait()
}

func multiGoroutineLogerExample() {
	if ll, err := lloger.Build(lloger.ERROR, "./varlog", "loop",
		0x80000, 10*(1<<20), 2*time.Second); err == nil {
		defer lloger.Close(ll)

		handlerpool(ll)
	}
}
```

###### 高并发时注意:

```go
//当写入日志并发量过高时, 除error级别外的日志, 都不会阻塞等待写入,会直接被丢弃

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
```

