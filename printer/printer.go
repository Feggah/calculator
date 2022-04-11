package printer

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Feggah/calculator/utils"
	"github.com/alexbrainman/printer"
)

const (
	timezone   = "America/Sao_Paulo"
	errPrinter = "Um erro inesperado ocorreu, entre em contato com um administrador."
)

type Printer struct {
	// Title of the document to be printed
	title string

	// Content of the document
	content string

	// Timestamp when the document was printed
	timestamp string

	// Current window configuration.
	// Used to manage the pop-up
	window fyne.Window
}

func NewPrinter(w fyne.Window, c string) *Printer {
	return &Printer{
		window:    w,
		content:   c,
		timestamp: getLocalTimestamp(),
	}
}

func (p *Printer) ShowPrinterPopUp() {
	input := widget.NewEntry()
	input.SetPlaceHolder("Insira o título...")

	var modal *widget.PopUp
	confirm := widget.NewButton("Confirmar", func() {
		p.title = input.Text
		modal.Hide()
		p.print()
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
		p.showErrorPopUp(err)
		return
	}

	prt, err := printer.Open(name)
	if err != nil {
		p.showErrorPopUp(err)
		return
	}

	err = prt.StartRawDocument(p.title)
	if err != nil {
		p.showErrorPopUp(err)
		return
	}

	content := "Titulo: " + p.title + "\n" + "Horario: " + p.timestamp + "\n\n" + parsedContent
	_, err = prt.Write([]byte(content))
	if err != nil {
		p.showErrorPopUp(err)
		return
	}

	err = prt.EndDocument()
	if err != nil {
		p.showErrorPopUp(err)
		return
	}
	err = prt.Close()
	if err != nil {
		p.showErrorPopUp(err)
		return
	}
}

func (p *Printer) showErrorPopUp(err error) {
	var modal *widget.PopUp
	ok := widget.NewButton("Ok", func() {
		modal.Hide()
	})
	ok.Importance = widget.HighImportance

	errMessage := &widget.Label{
		TextStyle: fyne.TextStyle{Monospace: true},
		Wrapping:  fyne.TextWrapWord,
		Text:      errPrinter,
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
		p.window.Canvas(),
	)
	modal.Resize(fyne.NewSize(300, 200))
	modal.Show()
}

func parseContent(content string) string {
	var parsedContent string
	for pos, char := range content {
		if content[pos] == '=' {
			parsedContent += string("\n---------\nTotal: ")
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
