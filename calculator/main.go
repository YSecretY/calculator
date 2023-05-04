package main

import "fyne.io/fyne/v2/app"

func main() {
	a := app.New()

	c := newCalculator()
	c.loadUI(a)

	a.Run()
}
