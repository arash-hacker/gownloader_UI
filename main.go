package main

import (
	"fmt"
	g "gownloader/gownload"
	"math/rand"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type calc struct {
	equation  string
	functions map[string]func()

	output  *widget.Label
	buttons map[string]*widget.Button
	window  fyne.Window
}

func (c *calc) display(newtext string) {
	c.equation = newtext
	c.output.SetText(newtext)
}

func (c *calc) character(char rune) {
	c.display(c.equation + string(char))
}

func (c *calc) digit(d int) {
	r := rune(d)
	r += '0'
	c.character(r)
}

func (c *calc) clear() {
	c.display("")
}

func (c *calc) addButton(text string, action func()) *widget.Button {
	button := widget.NewButton(text, action)
	c.buttons[text] = button

	return button
}

func (c *calc) digitButton(number int) *widget.Button {
	str := fmt.Sprintf("%d", number)
	action := func() {
		c.digit(number)
	}
	c.functions[str] = action

	return c.addButton(str, action)
}

func (c *calc) charButton(char rune) *widget.Button {
	action := func() {
		c.character(char)
	}
	c.functions[string(char)] = action

	return c.addButton(string(char), action)
}

func (c *calc) typedRune(r rune) {
	if r == '=' {

		return
	} else if r == 'c' {
		c.clear()
		return
	}

	action := c.functions[string(r)]
	if action != nil {
		action()
	}
}

func (c *calc) typedKey(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {

		return
	}
}

func (c *calc) loadUI(app fyne.App) {
	c.output = widget.NewLabel("")
	c.output.Alignment = fyne.TextAlignTrailing
	c.output.TextStyle.Monospace = true
	equals := c.addButton("=", func() {

	})
	equals.Style = widget.PrimaryButton

	c.window = app.NewWindow("Calc")
	//c.window.SetIcon(icon.CalculatorBitmap)
	c.window.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		c.output,
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			c.charButton('+'),
			c.charButton('-'),
			c.charButton('*'),
			c.charButton('/')),
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			c.digitButton(7),
			c.digitButton(8),
			c.digitButton(9),
			c.addButton("C", func() {
				c.clear()
			})),
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			c.digitButton(4),
			c.digitButton(5),
			c.digitButton(6),
			c.charButton('(')),
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			c.digitButton(1),
			c.digitButton(2),
			c.digitButton(3),
			c.charButton(')')),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				c.digitButton(0),
				c.charButton('.')),
			equals)),
	)

	c.window.Canvas().SetOnTypedRune(c.typedRune)
	c.window.Canvas().SetOnTypedKey(c.typedKey)
	c.window.Show()
}

func newCalculator() *calc {
	c := &calc{}
	c.functions = make(map[string]func())
	c.buttons = make(map[string]*widget.Button)

	return c
}

// Show loads a calculator example window for the specified app context
func Show(app fyne.App) {
	c := newCalculator()
	c.loadUI(app)
}

func main() {
	a := app.New()
	w := a.NewWindow("Example")

	grid := fyne.NewContainerWithLayout(layout.NewGridLayout(2))
	group2 := widget.NewGroup("Widgets")
	group3 := widget.NewGroup("Widgets")

	entry := widget.NewMultiLineEntry()
	stop := widget.NewButton("download", func() {
		group := widget.NewGroup((entry.Text)[:10])
		//entry.Text=""
		p := widget.NewProgressBar()
		b := widget.NewButton("start/stop", func() {})
		group.Append(p)
		group.Append(b)
		p.SetValue(rand.Float64())
		group3.Append(group)
		d := g.New(entry.Text)
		d.Init(4)
		b.Text = fmt.Sprintf("%d", d.GetSize())
		go d.StartAll2()
		tkr := time.NewTicker(100 * time.Millisecond)
		go func() {
			for {
				select {
				case <-tkr.C: 
					p.SetValue(d.Check()/float64(d.GetSize()))
				} 
			}
		}()

	})
	
	resume := widget.NewButton("stop", func() {})

	group2.Append(stop)
	group2.Append(resume)
	group3.Resize(fyne.NewSize(50, 800))
	entry.Resize(fyne.NewSize(50, 800))
	group3.Append(entry)

	grid.AddObject(group2)
	grid.AddObject(group3)

	w.SetContent(grid)
	w.ShowAndRun()
}
