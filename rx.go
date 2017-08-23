package rxgo

import (
	"reflect"
)

// Stream ...
type Stream interface {
	Next() (interface{}, bool)
}

type stream chan interface{}

// New ...
func New(from interface{}) Stream {
	if from == nil {
		return nil
	}
	refFrom := reflect.ValueOf(from)

	switch refFrom.Type().Kind() {

	case reflect.Chan:
		if refFrom.Type().ChanDir() == reflect.SendDir {
			return nil
		}

		ch := make(chan interface{})
		go func() {
			for {
				v, ok := refFrom.TryRecv()
				if !ok {
					break
				}
				ch <- v.Interface()
			}
			close(ch)
		}()
		return stream(ch)

	default:
		return nil
	}
}

// Next ...
func (s stream) Next() (interface{}, bool) {
	value, ok := <-s
	return value, ok
}
