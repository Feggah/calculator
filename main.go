package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/Feggah/calculator/calculator"
)

func main() {
	app := app.New()

	c := calculator.NewCalculator()
	c.LoadUserInterface(app)
	app.Run()
}
