package manifest

import (
	"fmt"
	"sync"
)

type Stream chan string

type Output struct {
	lock     sync.Mutex
	prefixes map[Stream]string
	streams  map[string]Stream
}

func NewOutput() Output {
	return Output{
		prefixes: make(map[Stream]string),
		streams:  make(map[string]Stream),
	}
}

func (o *Output) Stream(prefix string) Stream {
	if s, ok := o.streams[prefix]; ok {
		return s
	}

	s := make(Stream)

	o.prefixes[s] = prefix
	o.streams[prefix] = s

	go o.watchStream(s)

	return s
}

func (o *Output) paddedPrefix(s Stream) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", o.widestPrefix()), o.prefixes[s])
}

func (o *Output) printLine(s Stream, line string) {
	o.lock.Lock()
	defer o.lock.Unlock()

	fmt.Printf("%s | %s\n", o.paddedPrefix(s), line)
}

func (o *Output) watchStream(s Stream) {
	for line := range s {
		o.printLine(s, line)
	}
}

func (o *Output) widestPrefix() (w int) {
	for _, p := range o.prefixes {
		if len(p) > w {
			w = len(p)
		}
	}

	return
}
