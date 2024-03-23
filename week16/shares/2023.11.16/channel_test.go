package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestChannel(t *testing.T) {
	c := make(chan struct{})
	close(c)
	assert.Panics(t, func() {
		// 往已经关闭的channel写数据，会 panic
		c <- struct{}{}
	})

	//c := getFromXXX()

	// data := <- c
}

func TestChannelBlocking(t *testing.T) {
	// 没有初始化，c == nil
	var c chan struct{}
	// 这两个都会导致阻塞

	go func() {
		<-c
		t.Log("111不会输出这一句")
	}()
	var b1 BigStruct
	go func() {
		var b2 BigStruct
		c <- struct{}{}
		t.Log("222不会输出这一句", b1, b2)
	}()
}

type BigStruct struct {
}

func Close[T io.Closer](t T) {
	t.Close()
}

//public class Abc<T extends Closer> {
//
//}
