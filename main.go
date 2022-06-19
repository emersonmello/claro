/*
Copyright © 2022  Emerson Mello

*/
package main

import (
	"github.com/emersonmello/claro/cmd"
)

func main() {
	cmd.Execute()
}

// func init() {
// 	c := make(chan os.Signal)
// 	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
// 	go func() {
// 		// handling CTRL + C signal
// 		<-c
// should I remove the output directory created by clone command?
// if yes, use ioutil.ReadDir() and os.RemoveAll()
// 		os.Exit(1)
// 	}()
// }
