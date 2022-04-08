package printer

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Feggah/calculator/utils"
)

const (
	timezone = "America/Sao_Paulo"
)

type Printer struct {
	title string

	content string

	timestamp string

	window fyne.Window
}

func NewPrinter(w fyne.Window, c string) *Printer {
	return &Printer{
		window:    w,
		content:   c,
		timestamp: getLocalTimestamp(),
	}
}

func (p *Printer) SetTitle() {
	input := widget.NewEntry()
	input.SetPlaceHolder("Insira o título...")

	confirm := widget.NewButton("Confirmar", func() {
		p.title = input.Text
		p.print()
	})
	confirm.Importance = widget.HighImportance

	var modal *widget.PopUp
	modal = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Título da impressão"),
			input,
			container.NewGridWithColumns(2,
				widget.NewButton("Cancelar", func() { modal.Hide() }),
				confirm,
			),
		),
		p.window.Canvas(),
	)
	modal.Show()
}

func (p *Printer) print() {
	parsedContent := parseContent(p.content)
	fmt.Println(p.title)
	fmt.Println(p.timestamp)
	fmt.Println(parsedContent)
}

func parseContent(content string) string {
	var parsedContent string
	for pos, char := range content {
		if content[pos] == '=' {
			parsedContent += string("\n-----\n")
		} else {
			parsedContent += string(char)
		}
		if pos+1 < len(content) && utils.ByteInSlice(content[pos+1], utils.Operators) {
			parsedContent += "\n"
		}
	}
	return parsedContent
}

func getLocalTimestamp() string {
	loc, _ := time.LoadLocation(timezone)
	return time.Now().In(loc).String()
}
