package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("allocating new bytes.Buffer")
		return new(bytes.Buffer)
	},
}

// TODO: create pool of bytes.Buffers which can be reused
func log(w io.Writer, debug string){

	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()

	b.WriteString(time.Now().Format("15:04:05"))
	b.WriteString(" : ")
	b.WriteString(debug)
	b.WriteString("\n")

	w.Write(b.Bytes())

	bufPool.Put(b)
}

func main(){
	log(os.Stdout, "debug-string1")
	log(os.Stdout, "debug-string-2")
}