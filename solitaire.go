package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/mattn/go-runewidth"
)

var st = tcell.StyleDefault
var black = st
var red = st.Foreground(tcell.ColorRed)

var qwerty = []rune("qw erty")
var asdfghhj = []rune("asdfghi")

var suits = [...]string{"♣", "♦", "♠", "♥"}
var ranks = [...]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
var colors = [...]tcell.Style{black, red}

type Tabletop struct {
	stock, talon []int
	tableau      [7][]int
	foundations  [4][]int
	hidden       [7]int
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
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

func (t *Tabletop) Deal() {
	t.stock = rand.Perm(52)
	t.hidden = [...]int{0, 1, 2, 3, 4, 5, 6}
	for i := 1; i <= 7; i++ {
		t.tableau[i-1] = t.stock[len(t.stock)-i:]
		t.stock = t.stock[:len(t.stock)-i]
	}
	t.foundations[0] = []int{26}
	t.foundations[2] = []int{4}
}

func (t Tabletop) DrawCard(s tcell.Screen, x, y int, c int) {
	var str string
	if c == 52 {
		str = " ? "
	} else {
		str = fmt.Sprintf("%-2s%s", ranks[c%13], suits[c%4])
	}
	emitStr(s, x, y, colors[c%2], str)
}

func (t *Tabletop) Draw(s tcell.Screen) {
	for i, r := range qwerty {
		s.SetContent(i+1, 1, r, []rune("     "), black)
	}
	if len(t.stock) > 0 {
		t.DrawCard(s, 0, 3, t.stock[len(t.stock)-1])
	}
	if len(t.talon) > 0 {
		t.DrawCard(s, 6, 3, t.talon[len(t.talon)-1])
	}
	for i := 0; i <= 3; i++ {
		if len(t.foundations[i]) > 0 {
			t.DrawCard(s, 6*(i+3), 3, t.foundations[i][len(t.foundations[i])-1])
		}
	}

	for i, r := range asdfghhj {
		s.SetContent(i+1, 6, r, []rune("     "), black)
	}
	for i := 0; i <= 6; i++ {
		for j := 0; j < t.hidden[i]; j++ {
			t.DrawCard(s, 6*i, 8+j, 52)
		}
		for j := t.hidden[i]; j < len(t.tableau[i]); j++ {
			t.DrawCard(s, 6*i, 8+j, t.tableau[i][j])
		}
	}
}

func main() {
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	encoding.Register()

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()

	quit := make(chan struct{})

	var t Tabletop

	t.Deal()

	t.Draw(s)

	s.Show()
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	<-quit

	s.Fini()
}
