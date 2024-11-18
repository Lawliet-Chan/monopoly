package main

import (
	"github.com/fzdwx/infinite"
	"github.com/fzdwx/infinite/components/input/confirm"
)

func main() {

	c := infinite.NewConfirm(
		confirm.WithDefaultYes(),
		confirm.WithDisplayHelp(),
	)
	c.Display()
}
