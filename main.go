package main

import (
	"prio/prio-cli/editor"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	editor := &editor.Editor{}
	err = editor.OpenTopLevel()
	if err != nil {
		panic(err)
	}
	err = editor.Edit()
	if err != nil {
		panic(err)
	}
	err = editor.Save()
	if err != nil {
		panic(err)
	}
}
