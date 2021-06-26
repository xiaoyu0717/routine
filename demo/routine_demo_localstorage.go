package main

import (
	"fmt"
	"github.com/haima/routine"
	"time"
)

func main() {
	routine.GetLocalStorage().Set("key","hello world")
	fmt.Println("name: ", routine.GetLocalStorage().Get("key"))

	// other goroutine cannot read it
	go func() {
		fmt.Println("name1: ", routine.GetLocalStorage().Get("key"))
	}()

	// but, the new goroutine could inherit/copy all local data from the current goroutine like this:
	routine.Go(func() {
		fmt.Println("name2: ", routine.GetLocalStorage().Get("key"))
	})

	// or, you could copy all local data manually
	ic := routine.BackupContext()
	go func() {
		routine.InheritContext(ic)
		fmt.Println("name3: ", routine.GetLocalStorage().Get("key"))
	}()

	time.Sleep(time.Second)
}
