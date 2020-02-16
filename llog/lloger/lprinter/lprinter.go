package lprinter

import (
	"fmt"
	"os"
)

//Printer ...
type Printer struct {
	core *os.File
}

//Println ...
func (p Printer) Println(prefix string, objs ...interface{}) (err error) {

	_, err = fmt.Fprintf(p.core, "%s %v\n", prefix, objs)
	return
}

//Printfl ...
func (p Printer) Printfl(prefix string, format string, objs ...interface{}) (err error) {
	_, err = fmt.Fprintf(p.core, fmt.Sprintf("%s %s\n", prefix, format), objs...)
	return
}

// Build ...
func Build(core *os.File) (p Printer) {
	p.core = core
	return
}

// Info ...
func Info(p Printer) (fi os.FileInfo, err error) {
	fi, err = p.core.Stat()
	return
}

// Close ...
func Close(p Printer) (err error) {
	err = p.core.Close()
	return
}
