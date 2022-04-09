package printer

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	// "github.com/alexbrainman/printer"

	"github.com/Feggah/calculator/utils"
	"github.com/alexbrainman/printer"
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

	var modal *widget.PopUp
	confirm := widget.NewButton("Confirmar", func() {
		p.title = "\t" + input.Text
		p.print()
		modal.Hide()
	})
	confirm.Importance = widget.HighImportance

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
	name, err := printer.Default()
	if err != nil {
		fmt.Println(err.Error())
	}

	prt, err := printer.Open(name)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = prt.StartDocument(p.title, "text")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = prt.StartPage()
	if err != nil {
		fmt.Println(err.Error())
	}

	content := p.title + "\n" + p.timestamp + "\n\n" + parsedContent
	_, err = prt.Write([]byte(content))
	if err != nil {
		fmt.Println(err.Error())
	}

	err = prt.EndPage()
	if err != nil {
		fmt.Println(err.Error())
	}
	err = prt.EndDocument()
	if err != nil {
		fmt.Println(err.Error())
	}
	err = prt.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
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
	t := time.Now().In(loc)
	return t.Format("02/01/2006 15:04")
}
