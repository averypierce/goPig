/*
Avery Pierce VanKirk 2017
Based on tcell mouse demo - github.com/gdamore/tcell/blob/master/_demos/mouse.go
*/

//Copyright 2015 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"unicode"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"

	"github.com/mattn/go-runewidth"
)

//EmitRune prints a single rune to specified coordinate
func EmitRune(s tcell.Screen, x, y int, style tcell.Style, c rune) {

	var comb []rune
	w := runewidth.RuneWidth(c)
	if w == 0 {
		comb = []rune{c}
		c = ' '
		w = 1
	}
	s.SetContent(x, y, c, comb, style)
}

//InputArea generates an area where you can enter and backspace text one rune at a time. Input is
//not stored in any data structure, just printed to screen.
//TODO: Add scrolling support with a textbuffer
func InputArea(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) func(key *tcell.EventKey) {

	//persistant storage for cursor location
	cx, cy := x1, y1

	return func(key *tcell.EventKey) {
		r := key.Rune()
		if key.Key() == tcell.KeyBackspace || key.Key() == tcell.KeyBackspace2 {

			if cx > x1 {
				cx--
				EmitRune(s, cx, cy, style, ' ')
			} else if cy > y1 {
				//set cursor to prev line
				cy--
				cx = x2
				EmitRune(s, cx, cy, style, ' ')
			}

		} else if unicode.IsLetter(r) {

			EmitRune(s, cx, cy, style, r)
			cx++
		}
		//check for textArea boundary
		if cx > x2 {
			cx = x1
			cy++
		}
	}
}

/*
textArea + bordered thing

and lets have the textbuffer be its own type, with its own logging methods
and then compositie that all together
*/

type textBox struct {
	x1, y1, x2, y2, cx, cy int
	style                  tcell.Style
}

//DrawPigLayout uses tcell mouse demo functions to draw default layout for pig
func DrawPigLayout(s tcell.Screen, c tcell.Style) {
	w, h := s.Size()

	qw := w / 2
	qh := h/2 - 3

	drawBox(s, 0, 0, qw, qh, c, ' ')
	drawBox(s, qw, 0, qw*2, qh, c, ' ')
	drawBox(s, 0, qh, qw, qh*2, c, ' ')
	drawBox(s, qw, qh, qw*2, qh*2, c, ' ')
	drawBox(s, 0, qh*2, qw*2, qh*2+6, c, ' ')

	EmitStr(s, 5, 0, c, " Box One ")
	EmitStr(s, qw+5, 0, c, " Box Two ")
	EmitStr(s, 5, qh, c, " Box Three ")
	EmitStr(s, qw+5, qh, c, " Box Four ")
	//emitRune(s, 1, qh*2+5, c, '>')

}

//EmitStr is part of tcell mouse demo. prints a string to specified coordinate
func EmitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

//drawBox is part of tcell mouse demo. generates an empty box at specified coordinates
func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, r rune) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}
	if y1 != y2 && x1 != x2 {
		// Only add corners if we need to
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetContent(col, row, r, nil, style)
		}
	}
}

//MouseDemoMain is old main() from tcell mouse demo.
//shows simple mouse and keyboard events.  Press ESC twice to exit.
func MouseDemoMain() {

	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	defStyle := tcell.StyleDefault.
		Background(tcell.ColorPurple).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.Clear()

	posfmt := "Mouse: %d, %d  "
	btnfmt := "Buttons: %s"
	keyfmt := "Keys: %s"
	white := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)

	mx, my := -1, -1
	ox, oy := -1, -1
	//bx, by := -1, -1
	w, h := s.Size()
	qw := w / 2
	qh := h/2 - 3

	lchar := '*'
	bstr := ""
	lks := ""
	ecnt := 0
	//drawBox(s, 0, qh*2, qw*2, qh*2+6, c, ' ')
	DrawPigLayout(s, white)
	mint := InputArea(s, 1, qh*2+1, qw*2-1, qh*2+6-1, white)
	//drawBox(s, qw, qh, qw*2, qh*2, c, ' ')
	quadrant4 := InputArea(s, qw+1, qh+1, qw*2-1, qh*2-1, white)

	for {

		EmitStr(s, 2, 2, white, "Press ESC twice to exit, C to clear.")
		EmitStr(s, 2, 3, white, fmt.Sprintf(posfmt, mx, my))
		EmitStr(s, 2, 4, white, fmt.Sprintf(btnfmt, bstr))
		EmitStr(s, 2, 5, white, fmt.Sprintf(keyfmt, lks))
		s.Show()
		bstr = ""
		ev := s.PollEvent()
		st := tcell.StyleDefault.Background(tcell.ColorRed)
		up := tcell.StyleDefault.
			Background(tcell.ColorBlue).
			Foreground(tcell.ColorBlack)
		w, h = s.Size()

		// always clear any old selection box
		/*if ox >= 0 && oy >= 0 && bx >= 0 {
			drawSelect(s, ox, oy, bx, by, false)
		}*/

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			s.SetContent(w-1, h-1, 'R', nil, st)
		case *tcell.EventKey:
			mint(ev)
			quadrant4(ev)
			//s.SetContent(w-2, h-2, ev.Rune(), nil, st)
			s.SetContent(w-1, h-1, 'K', nil, st)
			if ev.Key() == tcell.KeyEscape {
				ecnt++
				if ecnt > 1 {
					s.Fini()
					os.Exit(0)
				}
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else {
				ecnt = 0
				/*if ev.Rune() == 'C' || ev.Rune() == 'c' {
					s.Clear()
				}*/
			}
			lks = ev.Name()
		case *tcell.EventMouse:
			x, y := ev.Position()
			button := ev.Buttons()
			for i := uint(0); i < 8; i++ {
				if int(button)&(1<<i) != 0 {
					bstr += fmt.Sprintf(" Button%d", i+1)
				}
			}
			if button&tcell.WheelUp != 0 {
				bstr += " WheelUp"
			}
			if button&tcell.WheelDown != 0 {
				bstr += " WheelDown"
			}
			if button&tcell.WheelLeft != 0 {
				bstr += " WheelLeft"
			}
			if button&tcell.WheelRight != 0 {
				bstr += " WheelRight"
			}
			// Only buttons, not wheel events
			button &= tcell.ButtonMask(0xff)
			ch := '*'

			if button != tcell.ButtonNone && ox < 0 {
				ox, oy = x, y
			}
			switch ev.Buttons() {
			case tcell.ButtonNone:
				if ox >= 0 {
					bg := tcell.Color((lchar - '0') * 2)
					drawBox(s, ox, oy, x, y,
						up.Background(bg),
						lchar)
					ox, oy = -1, -1
					//bx, by = -1, -1
				}
			case tcell.Button1:
				ch = '1'
			case tcell.Button2:
				ch = '2'
			case tcell.Button3:
				ch = '3'
			case tcell.Button4:
				ch = '4'
			case tcell.Button5:
				ch = '5'
			case tcell.Button6:
				ch = '6'
			case tcell.Button7:
				ch = '7'
			case tcell.Button8:
				ch = '8'
			default:
				ch = '*'

			}
			/*if button != tcell.ButtonNone {
				bx, by = x, y
			}*/
			lchar = ch
			//s.SetContent(w-1, h-1, 'M', nil, st)
			mx, my = x, y
		default:
			s.SetContent(w-1, h-1, 'X', nil, st)
		}

		/*
			if ox >= 0 && bx >= 0 {
				drawSelect(s, ox, oy, bx, by, true)
			}
		*/
	}
}
