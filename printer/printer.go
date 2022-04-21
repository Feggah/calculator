package printer

import (
	"os"
	"os/exec"
	"time"

	_ "time/tzdata"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Feggah/calculator/utils"
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

	// Current window configuration.
	// Used to manage the pop-up
	window fyne.Window
}

func NewPrinter(w fyne.Window, c string) *Printer {
	return &Printer{
		window:  w,
		content: c,
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
	content := "Titulo: " + p.title + "\n" + "Horario: " + getLocalTimestamp() + "\n\n" + parseContent(p.content)

	file, err := os.CreateTemp("", p.title)
	if err != nil {
		p.showErrorPopUp(err)
		return
	}

	defer func() {
		file.Close()
		if err := os.Remove(file.Name()); err != nil {
			p.showErrorPopUp(err)
			return
		}
	}()

	if _, err := file.WriteString(content); err != nil {
		p.showErrorPopUp(err)
		return
	}

	cmd := exec.Command("cmd.exe", "/C", "notepad", "/p", file.Name())
	if err := cmd.Run(); err != nil {
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
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return errPrinter + " 'Error when getting the local timezone location'"
	}
	t := time.Now().In(loc)
	return t.Format("02/01/2006 15:04")
}
