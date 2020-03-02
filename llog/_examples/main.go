package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/llog-master/llog/lloger"
)

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

func consolelogerExample() {
	ll := lloger.BuildConsoleLoger(lloger.DEBUG)
	for i := uint(0); i < 0xffff; i++ {
		ll.Debugfl("%#v", "hello, world")
		ll.Debugln("hello, world")
	}

}

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

func main() {
	fmt.Println(time.Now())
	multiGoroutineLogerExample()
	fmt.Println(time.Now())
}
