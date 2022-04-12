package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/Feggah/calculator/calculator"
	"github.com/Feggah/calculator/utils"
)

func main() {
	app := app.New()
	app.SetIcon(utils.ResourceIconPng)

	c := calculator.NewCalculator()
	c.LoadUserInterface(app)
	app.Run()
}
