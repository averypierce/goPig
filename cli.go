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
	"bytes"
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

type textBox struct {
	x1, y1, x2, y2, cx, cy int
	title                  string
	content                string
	s                      tcell.Screen
	style                  tcell.Style
	buffer                 bytes.Buffer
}

//All this logic needs to run on runes, NOT tcell events if we want resizing to work. or...no it doesnt. because text will be stored..hmms
func (tb *textBox) drawBoarder() {
	tb.cx = tb.x1 + 1
	tb.cy = tb.y1 + 1
	drawBox(tb.s, tb.x1, tb.y1, tb.x2, tb.y2, tb.style, ' ')
	EmitStr(tb.s, tb.x1+5, tb.y1, tb.style, tb.title)
}

func (tb *textBox) input(key *tcell.EventKey) {
	r := key.Rune()
	if key.Key() == tcell.KeyBackspace || key.Key() == tcell.KeyBackspace2 {

		if tb.cx > tb.x1+1 {
			tb.cx--
			EmitRune(tb.s, tb.cx, tb.cy, tb.style, ' ')
		} else if tb.cy > tb.y1+1 {
			//set cursor to prev line
			tb.cy--
			tb.cx = tb.x2 - 1
			EmitRune(tb.s, tb.cx, tb.cy, tb.style, ' ')
		}

	} else if !unicode.IsControl(r) {

		EmitRune(tb.s, tb.cx, tb.cy, tb.style, r)
		tb.cx++
		tb.buffer.WriteRune(r)
	}
	//check for textArea boundary
	if tb.cx > tb.x2-1 {
		tb.cx = tb.x1 + 1
		tb.cy++
	}
}

func (tb *textBox) redraw(x1, y1, x2, y2 int) {
	tb.x1, tb.y1, tb.x2, tb.y2 = x1, y1, x2, y2
	tb.drawBoarder()

	buf := tb.buffer.String()

	for _, r := range buf {
		EmitRune(tb.s, tb.cx, tb.cy, tb.style, r)
		tb.cx++

		if tb.cx > tb.x2-1 {
			tb.cx = tb.x1 + 1
			tb.cy++
		}
	}

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
	s.Clear()

	posfmt := "Mouse: %d, %d  "
	btnfmt := "Buttons: %s"
	keyfmt := "Keys: %s"
	white := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)

	mx, my := -1, -1
	w, h := s.Size()
	qw := w / 2
	qh := h/2 - 3

	bstr := ""
	lks := ""
	ecnt := 0

	q1 := textBox{x1: 0, y1: 0, x2: qw, y2: qh, title: " Box One ", content: "", s: s, style: white}
	q2 := textBox{x1: qw, y1: 0, x2: qw * 2, y2: qh, title: " Box Two ", content: "", s: s, style: white}
	q3 := textBox{x1: 0, y1: qh, x2: qw, y2: qh * 2, title: " Box Three ", content: "", s: s, style: white}
	q4 := textBox{x1: qw, y1: qh, x2: qw * 2, y2: qh * 2, title: " Box Four ", content: "", s: s, style: white}
	cmd := textBox{x1: 0, y1: qh*2 + 1, x2: qw * 2, y2: qh*2 + 6, title: " Commands ", content: "", s: s, style: white}
	q1.drawBoarder()
	q2.drawBoarder()
	q3.drawBoarder()
	q4.drawBoarder()
	cmd.drawBoarder()

	for {

		EmitStr(s, 2, 2, white, "Press ESC twice to exit, C to clear.")
		EmitStr(s, 2, 3, white, fmt.Sprintf(posfmt, mx, my))
		EmitStr(s, 2, 4, white, fmt.Sprintf(btnfmt, bstr))
		EmitStr(s, 2, 5, white, fmt.Sprintf(keyfmt, lks))
		s.Show()
		bstr = ""
		ev := s.PollEvent()
		st := tcell.StyleDefault.Background(tcell.ColorRed)

		w, h = s.Size()
		qw = w / 2
		qh = h/2 - 3

		switch ev := ev.(type) {
		case *tcell.EventResize:
			q1.redraw(0, 0, qw, qh)
			q2.redraw(qw, 0, qw*2, qh)
			q3.redraw(0, qh, qw, qh*2)
			q4.redraw(qw, qh, qw*2, qh*2)
			cmd.redraw(0, qh*2, qw*2, qh*2+6)

			s.Sync()
		case *tcell.EventKey:

			q3.input(ev)
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
			}
			lks = ev.Name()

		default:
			s.SetContent(w-1, h-1, 'X', nil, st)
		}
	}
}
