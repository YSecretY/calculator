package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
)

// calc is a struct that implements calculator
type calc struct {
	output  *widget.Label
	res     string
	base    int // current base of the number
	prev    string
	lastOp  rune
	buttons map[string]*widget.Button
	window  fyne.Window
}

// newCalculator returns a pointer to calculator struct,
// it implements map with capacity equal to amount of buttons,
// output as a widget.NewLabel, value of previous number in output
// and last operation between numbers, it is essential to be +
func newCalculator() *calc {
	return &calc{
		buttons: make(map[string]*widget.Button, 24),
		output:  widget.NewLabel(""),
		base:    10,
		prev:    "0",
		lastOp:  '+',
	}
}

// display takes text from the res variable and
// shows it in output label
func (c *calc) display(text string) {
	c.res = text
	c.output.SetText(text)
}

// character adds char to res and displays it
func (c *calc) character(char rune) {
	c.display(c.res + string(char))
}

// operation is a func of every operation button
// as +, -, * or /. It works based on operation and sends
// what to do to evaluate func
func (c *calc) operation(op rune) {
	switch op {
	case '+':
		c.evaluate()
		c.lastOp = '+'
		prev, err := c.convertInt(c.output.Text, c.base, 10)
		if err != nil {
			log.Println("cannot convert number:", err)
		}
		c.prev = prev
	case '-':
		c.evaluate()
		c.lastOp = '-'
		prev, err := c.convertInt(c.output.Text, c.base, 10)
		if err != nil {
			log.Println("cannot convert number:", err)
		}
		c.prev = prev
	case '*':
		c.evaluate()
		c.lastOp = '*'
		prev, err := c.convertInt(c.output.Text, c.base, 10)
		if err != nil {
			log.Println("cannot convert number:", err)
		}
		c.prev = prev
	case '/':
		c.evaluate()
		c.lastOp = '/'
		prev, err := c.convertInt(c.output.Text, c.base, 10)
		if err != nil {
			log.Println("cannot convert number:", err)
		}
		c.prev = prev
	}
}

// digit adds digit to res and displays it
func (c *calc) digit(d int) {
	c.character(rune(d) + '0')
}

// clear clears the output, and resets basic variables
// to their base form
func (c *calc) clear() {
	c.display("")
	c.prev = "0"
	c.lastOp = '+'
}

// backspace deletes last symbol from the output
func (c *calc) backspace() {
	if len(c.res) == 0 {
		return
	} else if c.res == "error" {
		c.clear()
		return
	}

	c.display(c.res[:len(c.res)-1])
}

// addButton adds button to the map
func (c *calc) addButton(text string, action func()) *widget.Button {
	button := widget.NewButton(text, action)
	c.buttons[text] = button
	return button
}

// digitButton adds button to the map
// (by sending it to addButton) and adds
// digit func when click on it
func (c *calc) digitButton(number int) *widget.Button {
	str := strconv.Itoa(number)
	return c.addButton(str, func() {
		c.digit(number)
	})
}

// charButton adds button to the map
// (by sending it to addButton) and adds
// character func when click on it
func (c *calc) charButton(char rune) *widget.Button {
	return c.addButton(string(char), func() {
		c.character(char)
	})
}

// operationButton adds button to the map
// (by sending it to addButton) and adds
// operation func when click on it
func (c *calc) operationButton(op rune) *widget.Button {
	return c.addButton(string(op), func() {
		c.operation(op)
	})
}

// evaluate calculates result of last operation
// and changes output text to the result
// e.g.: 0.0 + 5 = 5
// 5 (current on the display) - 2 and if
// user clicks + here, evaluate still counts 5 - 2 = 3 at first
// and after that adds to 3 next number
func (c *calc) evaluate() {
	prev, err := strconv.ParseInt(c.prev, 10, 64)
	if err != nil {
		log.Println("cannot parse prev int:", err)
		c.output.SetText("error: cannot convert this number")
		return
	}

	var cur int64
	if c.output.Text != "" {
		cur, err = strconv.ParseInt(c.output.Text, c.base, 64)
		if err != nil {
			log.Println("cannot parse cur int:", err)
			c.output.SetText("error: cannot convert this number")
			return
		}
	}

	var res int64
	switch c.lastOp {
	case '+':
		res = prev + cur
		resInBase, err := c.convertInt(strconv.FormatInt(res, 10), 10, c.base)
		if err != nil {
			log.Println("unable to convert this number:", err)
		}
		c.clear()
		c.output.SetText(resInBase)
	case '-':
		res = prev - cur
		resInBase, err := c.convertInt(strconv.FormatInt(res, 10), 10, c.base)
		if err != nil {
			log.Println("unable to convert this number:", err)
		}
		c.clear()
		c.output.SetText(resInBase)
	case '*':
		res = prev * cur
		resInBase, err := c.convertInt(strconv.FormatInt(res, 10), 10, c.base)
		if err != nil {
			log.Println("unable to convert this number:", err)
		}
		c.clear()
		c.output.SetText(resInBase)
	case '/':
		res = prev / cur
		resInBase, err := c.convertInt(strconv.FormatInt(res, 10), 10, c.base)
		if err != nil {
			log.Println("unable to convert this number:", err)
		}
		c.clear()
		c.output.SetText(resInBase)
	}
}

// loadUI runs main windows of the program
// (calculator with buttons and output)
func (c *calc) loadUI(app fyne.App) {
	c.window = app.NewWindow("Calculator")
	c.window.Resize(fyne.NewSize(800, 600))
	c.addMenu()
	c.output.TextStyle.Monospace = true

	digits := []int{7, 8, 9, 4, 5, 6, 1, 2, 3, 0}
	for _, num := range digits {
		c.digitButton(num)
	}

	chars := []rune{'a', 'b', 'c', 'd', 'e', 'f'}
	for _, char := range chars {
		c.charButton(char)
	}

	operations := []rune{'+', '-', '*', '/'}
	for _, operation := range operations {
		c.operationButton(operation)
	}

	c.buttons["backspace"] = widget.NewButton("backspace", c.backspace)
	equals := c.addButton("=", c.evaluate)

	c.window.SetContent(container.NewGridWithRows(6,
		container.NewGridWithColumns(2, c.output, widget.NewButton("CE", c.clear)),
		container.NewGridWithColumns(6, c.buttons["a"], c.buttons["b"], c.buttons["c"], c.buttons["d"], c.buttons["e"], c.buttons["f"]),
		container.NewGridWithColumns(4, c.buttons["7"], c.buttons["8"], c.buttons["9"], c.buttons["backspace"]),
		container.NewGridWithColumns(4, c.buttons["4"], c.buttons["5"], c.buttons["6"], c.buttons["+"]),
		container.NewGridWithColumns(4, c.buttons["1"], c.buttons["2"], c.buttons["3"], c.buttons["-"]),
		container.NewGridWithColumns(4, equals, c.buttons["0"], c.buttons["/"], c.buttons["*"]),
	))
	c.window.Show()
}

// addMenu adds convert main menu in the top left corner
// it allows to convert result on the display to any base system
// and keep count there
func (c *calc) addMenu() {
	toBase := fyne.NewMenuItem("toBase", nil)
	childMenuItems := make([]*fyne.MenuItem, 15)
	prevChecked := -1
	for i := 0; i < 15; i++ {
		base := i + 2
		childMenuItems[i] = fyne.NewMenuItem(fmt.Sprintf("%d", base), func() {
			res, err := c.convertInt(c.output.Text, c.base, base)
			if err != nil {
				res = "error: cannot convert this number"
			}
			c.base = base
			c.output.SetText(res)
			childMenuItems[base-2].Checked = true
			if prevChecked != -1 {
				childMenuItems[prevChecked].Checked = false
			}
			prevChecked = base - 2
		})
	}
	toBase.ChildMenu = fyne.NewMenu("", childMenuItems...)

	convertMenu := fyne.NewMenu("Convert", toBase)
	mainMenu := fyne.NewMainMenu(convertMenu)
	c.window.SetMainMenu(mainMenu)
}

// convertInt takes string value of int and convert it from one base system to another
func (c *calc) convertInt(val string, fromBase int, toBase int) (string, error) {
	intPart, err := strconv.ParseInt(val, fromBase, 64)
	if err != nil {
		return "", err
	}
	res := strconv.FormatInt(intPart, toBase)
	return res, nil
}
