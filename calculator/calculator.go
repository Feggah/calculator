package calculator

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Feggah/calculator/printer"
	"github.com/Feggah/calculator/utils"
	"github.com/Knetic/govaluate"
)

const (
	errEquationExpressionInvalid = "Expressão matemática inválida. Tente novamente"
	title                        = "Calculadora"
)

type Calculator struct {
	// Window properties
	window fyne.Window

	// Calculator buttons
	buttons map[string]*widget.Button

	// Current equation calculation
	equation string

	// Output shown in the interface
	output *widget.Label

	// Controls the scrollable position
	// of the Output field
	outputPosition *container.Scroll

	// Last successfull equation
	lastEquation string
}

func NewCalculator() *Calculator {
	return &Calculator{
		buttons: make(map[string]*widget.Button, 19),
	}
}

func (c *Calculator) LoadUserInterface(app fyne.App) {
	c.configureOutput()
	c.outputPosition = container.NewVScroll(c.output)
	c.configureWindow(app)
}

func (c *Calculator) newButton(content string, callback func()) *widget.Button {
	b := widget.NewButton(content, callback)
	c.buttons[content] = b
	return b
}

func (c *Calculator) addNumericButton(number int) *widget.Button {
	callback := func() { c.numericCallback(number) }
	return c.newButton(strconv.Itoa(number), callback)
}

func (c *Calculator) numericCallback(number int) {
	c.removeExpressionInvalidMessage()
	c.display(c.equation + strconv.Itoa(number))
}

func (c *Calculator) addCharacterButton(char rune) *widget.Button {
	callback := func() { c.charCallback(char) }
	return c.newButton(string(char), callback)
}

func (c *Calculator) charCallback(char rune) {
	c.removeExpressionInvalidMessage()
	if utils.ByteInSlice(byte(char), utils.Operators) {
		if len(c.output.Text) == 0 {
			return
		}

		lastChar := c.output.Text[len(c.output.Text)-1]
		if byte(char) == lastChar {
			c.backspace()
		} else if utils.ByteInSlice(lastChar, utils.Operators) {
			c.backspace()
		}
	}
	c.display(c.equation + string(char))
}

func (c *Calculator) display(content string) {
	c.equation = content
	c.output.SetText(content)
	c.outputPosition.ScrollToBottom()
}

func (c *Calculator) clear() {
	c.lastEquation = ""
	c.display("")
}

func (c *Calculator) removeExpressionInvalidMessage() {
	if strings.Contains(c.output.Text, errEquationExpressionInvalid) {
		c.clear()
	}
}

func (c *Calculator) evaluate() {
	if strings.Contains(c.output.Text, errEquationExpressionInvalid) {
		c.display(errEquationExpressionInvalid)
		return
	}

	expression, err := govaluate.NewEvaluableExpression(c.output.Text)
	if err != nil {
		c.display(errEquationExpressionInvalid)
		return
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		c.display(errEquationExpressionInvalid)
		return
	}

	value, ok := result.(float64)
	if !ok {
		c.display(errEquationExpressionInvalid)
		return
	}

	output := strconv.FormatFloat(value, 'f', 2, 64)
	c.lastEquation = c.output.Text + "=" + output
	c.display(output)
}

func (c *Calculator) onTypedChar(char rune) {
	if string(char) == "," {
		char = '.'
	}
	if button, ok := c.buttons[string(char)]; ok {
		button.OnTapped()
	}
}

func (c *Calculator) onTypedKey(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter {
		if !c.isEquation() {
			// We don't need to evaluate if there isn't a equation in
			// the calculator display
			return
		}
		c.removeLastOperator()
		c.evaluate()
	} else if key.Name == fyne.KeyBackspace {
		c.backspace()
	}
}

func (c *Calculator) isEquation() bool {
	for _, v := range c.output.Text {
		if utils.ByteInSlice(byte(v), utils.Operators) {
			return true
		}
	}
	return false
}

func (c *Calculator) removeLastOperator() {
	if utils.ByteInSlice(c.output.Text[len(c.output.Text)-1], utils.Operators) {
		c.backspace()
	}
}

func (c *Calculator) backspace() {
	if len(c.equation) == 0 {
		return
	} else if c.equation == errEquationExpressionInvalid {
		c.clear()
		return
	}

	c.display(c.equation[:len(c.equation)-1])
}

func (c *Calculator) configureOutput() {
	c.output = &widget.Label{Alignment: fyne.TextAlignTrailing}
	c.output.TextStyle.Monospace = true
	c.output.Wrapping = fyne.TextWrapBreak
}

func (c *Calculator) configureWindow(app fyne.App) {
	equals := c.newButton("=", c.evaluate)
	equals.Importance = widget.HighImportance

	backspace := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), c.backspace)
	c.buttons["backspace"] = backspace

	c.window = app.NewWindow(title)
	c.window.SetContent(
		container.NewGridWithColumns(1,
			c.outputPosition,
			container.NewGridWithColumns(4,
				c.newButton("C", c.clear),
				c.addCharacterButton('%'),
				c.addCharacterButton('/'),
				backspace,
			),
			container.NewGridWithColumns(4,
				c.addNumericButton(7),
				c.addNumericButton(8),
				c.addNumericButton(9),
				c.addCharacterButton('*'),
			),
			container.NewGridWithColumns(4,
				c.addNumericButton(4),
				c.addNumericButton(5),
				c.addNumericButton(6),
				c.addCharacterButton('-'),
			),
			container.NewGridWithColumns(4,
				c.addNumericButton(1),
				c.addNumericButton(2),
				c.addNumericButton(3),
				c.addCharacterButton('+'),
			),
			container.NewGridWithColumns(2,
				container.NewGridWithColumns(2,
					c.addNumericButton(0),
					c.addCharacterButton('.')),
				equals,
			),
		),
	)

	ctrlP := desktop.CustomShortcut{KeyName: fyne.KeyP, Modifier: desktop.ControlModifier}
	c.window.Canvas().AddShortcut(&ctrlP, func(shortcut fyne.Shortcut) {
		c.printCallback()
	})
	c.window.Canvas().SetOnTypedRune(c.onTypedChar)
	c.window.Canvas().SetOnTypedKey(c.onTypedKey)
	c.window.Resize(fyne.NewSize(450, 500))
	c.configureMenu()
	c.window.Show()
}

func (c *Calculator) configureMenu() {
	item := fyne.NewMenuItem("Imprimir (CTRL + P)", c.printCallback)
	menu := fyne.NewMenu("Arquivo", item)
	main := fyne.NewMainMenu(menu)
	c.window.SetMainMenu(main)
}

func (c *Calculator) printCallback() {
	p, err := printer.NewPrinter(c.window, c.lastEquation)
	if err != nil {
		utils.ShowErrorPopUp(c.window, err)
		return
	}
	p.ShowPrinterPopUp()
}
