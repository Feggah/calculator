package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	Operators     = []byte{'+', '-', '*', '/', '%'}
	ErrUnexpected = "Um erro inesperado ocorreu, entre em contato com um administrador."
)

func ByteInSlice(char byte, list []byte) bool {
	for _, b := range list {
		if b == char {
			return true
		}
	}
	return false
}

func ShowErrorPopUp(window fyne.Window, err error) {
	var modal *widget.PopUp
	ok := widget.NewButton("Ok", func() {
		modal.Hide()
	})
	ok.Importance = widget.HighImportance

	errMessage := &widget.Label{
		TextStyle: fyne.TextStyle{Monospace: true},
		Wrapping:  fyne.TextWrapWord,
		Text:      ErrUnexpected,
	}

	errDetailed := &widget.Label{
		Text:      "Causa do erro: " + err.Error(),
		Wrapping:  fyne.TextWrapWord,
		TextStyle: fyne.TextStyle{Italic: true},
	}

	errTitle := &widget.Label{
		Text:      "Erro",
		TextStyle: fyne.TextStyle{Bold: true},
	}

	modal = widget.NewModalPopUp(
		container.NewVBox(
			errTitle,
			errMessage,
			errDetailed,
			ok,
		),
		window.Canvas(),
	)
	modal.Resize(fyne.NewSize(300, 200))
	modal.Show()
}
