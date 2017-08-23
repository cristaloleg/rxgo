package rxgo_test

import (
	"testing"

	"github.com/cristaloleg/rxgo"
)

func TestStream(t *testing.T) {
	tests := []struct {
		in interface{}
	}{
		{make(chan interface{})},
		{make(chan int)},
		{make(<-chan interface{})},
		{make(<-chan int)},
		{make(<-chan chan int)},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if rxgo.New(tt.in) == nil {
				t.Errorf("expected not nil `%v`", tt.in)
			}
		})
	}
}

func TestStreamNil(t *testing.T) {
	tests := []struct {
		in interface{}
	}{
		{nil},
		{"str"},
		{42},
		{make([]int, 0)},
		{make(chan<- int)},
		{make(chan<- chan int)},
		{func() {}},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := rxgo.New(tt.in); got != nil {
				t.Errorf("expected nil `%v` got `%v`", tt.in, got)
			}
		})
	}
}

func TestStreamNext(t *testing.T) {
	tests := []struct {
		data interface{}
	}{
		{10},
		{"string"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			ch := make(chan interface{})

			s := rxgo.New(ch)
			if s == nil {
				t.Errorf("expected not nil")
			}

			ch <- tt.data

			value, ok := s.Next()
			if !ok {
				t.Errorf("expected true")
			}

			if value != tt.data {
				t.Errorf("got %v want %v", value, tt.data)
			}

			_, ok = s.Next()
			if ok {
				t.Error("expected to be closed")
			}
		})
	}
}
