package printer

import (
	"errors"
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
	timezone = "America/Sao_Paulo"
	errEmpty = "não existe nenhuma conta registrada. Tente novamente depois de efetuar uma conta"
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

func NewPrinter(w fyne.Window, c string) (*Printer, error) {
	if c == "" {
		return nil, errors.New(errEmpty)
	}
	return &Printer{
		window:  w,
		content: c,
	}, nil
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
		utils.ShowErrorPopUp(p.window, err)
		return
	}

	defer func() {
		file.Close()
		if err := os.Remove(file.Name()); err != nil {
			utils.ShowErrorPopUp(p.window, err)
			return
		}
	}()

	if _, err := file.WriteString(content); err != nil {
		utils.ShowErrorPopUp(p.window, err)
		return
	}

	cmd := exec.Command("cmd.exe", "/C", "notepad", "/p", file.Name())
	if err := cmd.Run(); err != nil {
		utils.ShowErrorPopUp(p.window, err)
		return
	}
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
		return utils.ErrUnexpected + " 'Error when getting the local timezone location'"
	}
	t := time.Now().In(loc)
	return t.Format("02/01/2006 15:04")
}
