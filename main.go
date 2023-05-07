package main

import (
	"os"

	"github.com/Mahmoud-Khaled-FS/suha/app"
)

func main() {
	app := app.New(os.Args[1:])
	app.Start()
}
