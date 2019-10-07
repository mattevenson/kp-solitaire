package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/mattn/go-runewidth"
)

var st = tcell.StyleDefault
var black = st
var red = st.Foreground(tcell.ColorRed)

var qwerty = []rune("qw erty")
var asdfghhj = []rune("asdfghj")

var suits = [...]string{"♣", "♦", "♠", "♥"}
var ranks = [...]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
var colors = [...]tcell.Style{black, red}

var foundationsRunes = map[rune]int{
	'e': 0,
	'r': 1,
	't': 2,
	'y': 3,
}

var tableauRunes = map[rune]int{
	'a': 0,
	's': 1,
	'd': 2,
	'f': 3,
	'g': 4,
	'h': 5,
	'j': 6,
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

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

func CopyAndAppend(i []int, vals ...int) []int {
	j := make([]int, len(i), len(i)+len(vals))
	copy(j, i)
	return append(j, vals...)
}

type Tabletop struct {
	stock, talon []int
	tableau      [7][]int
	foundations  [4][]int
	hidden       [7]int
	key          rune
}

func (t *Tabletop) IsWon() bool {
	for _, foundation := range t.foundations {
		if len(foundation) < 13 {
			return false
		}
	}
	return true
}

func (t *Tabletop) Deal() {
	rand.Seed(time.Now().UnixNano())
	t.stock = rand.Perm(52)
	t.hidden = [...]int{0, 1, 2, 3, 4, 5, 6}
	for i := 1; i <= 7; i++ {
		t.tableau[i-1] = t.stock[len(t.stock)-i:]
		t.stock = t.stock[:len(t.stock)-i]
	}
	t.key = 'z'
}

func (t *Tabletop) CanPlayInTableau(c, i int) bool {
	if len(t.tableau[i]) == 0 {
		return c%13 == 12
	}
	c2 := t.tableau[i][len(t.tableau[i])-1]
	return c%13 == (c2-1)%13 && c%2 != c2%2
}

func (t *Tabletop) CanPlayInFoundations(c, i int) bool {
	if len(t.foundations[i]) == 0 {
		return c%13 == 0
	}
	c2 := t.foundations[i][len(t.foundations[i])-1]
	return c%13 == (c2+1)%13 && c%4 == c2%4
}

func (t *Tabletop) SelectKey(s tcell.Screen, key rune) {
	if t.key == 'z' {
		t.key = key
	} else {
		switch t.key {
		case 'q':
			if len(t.stock) == 0 {
				break
			}
			c := t.stock[len(t.stock)-1]
			switch key {
			case 'w':
				t.talon = append(t.talon, c)
				t.stock = t.stock[:len(t.stock)-1]
			case 'e', 'r', 't', 'y':
				i := foundationsRunes[key]
				if t.CanPlayInFoundations(c, i) {
					t.foundations[i] = copyAndAppend(t.foundations[i], c)
					t.stock = t.stock[:len(t.stock)-1]
				}
			case 'a', 's', 'd', 'f', 'g', 'h', 'j':
				i := tableauRunes[key]
				if t.CanPlayInTableau(c, i) {
					t.tableau[i] = copyAndAppend(t.tableau[i], c)
					t.stock = t.stock[:len(t.stock)-1]
				}
			}
		case 'w':
			if len(t.talon) == 0 {
				break
			}
			c := t.talon[len(t.talon)-1]
			switch key {
			case 'e', 'r', 't', 'y':
				i := foundationsRunes[key]
				if t.CanPlayInFoundations(c, i) {
					t.foundations[i] = copyAndAppend(t.foundations[i], c)
					t.talon = t.talon[:len(t.talon)-1]
				}
			case 'a', 's', 'd', 'f', 'g', 'h', 'j':
				i := tableauRunes[key]
				if t.CanPlayInTableau(c, i) {
					t.tableau[i] = copyAndAppend(t.tableau[i], c)
					t.talon = t.talon[:len(t.talon)-1]
				}
			}
		case 'e', 'r', 't', 'y':
			i := foundationsRunes[t.key]
			if len(t.foundations[i]) == 0 {
				break
			}
			c := t.foundations[i][len(t.foundations[i])-1]
			switch key {
			case 'w':
				t.talon = append(t.talon, c)
				t.foundations[i] = t.foundations[i][:len(t.foundations[i])-1]
			case 'e', 'r', 't', 'y':
				j := foundationsRunes[key]
				if t.CanPlayInFoundations(c, j) {
					t.foundations[j] = copyAndAppend(t.foundations[j], c)
					t.foundations[i] = t.foundations[i][:len(t.foundations[i])-1]
				}
			case 'a', 's', 'd', 'f', 'g', 'h', 'j':
				j := tableauRunes[key]
				if t.CanPlayInTableau(c, j) {
					t.tableau[j] = copyAndAppend(t.tableau[j], c)
					t.foundations[i] = t.foundations[i][:len(t.foundations[i])-1]
				}
			}
		case 'a', 's', 'd', 'f', 'g', 'h', 'j':
			i := tableauRunes[t.key]
			if len(t.tableau[i]) == 0 {
				break
			}
			c := t.tableau[i][len(t.tableau[i])-1]
			switch key {
			case 'e', 'r', 't', 'y':
				j := foundationsRunes[key]
				if t.CanPlayInFoundations(c, j) {
					t.foundations[j] = copyAndAppend(t.foundations[j], c)
					t.tableau[i] = t.tableau[i][:len(t.tableau[i])-1]
				}
			case 'a', 's', 'd', 'f', 'g', 'h', 'j':
				j := tableauRunes[key]
				if t.CanPlayInTableau(c, j) {
					t.tableau[j] = copyAndAppend(t.tableau[j], c)
					t.tableau[i] = t.tableau[i][:len(t.tableau[i])-1]
				}
			}
			if len(t.tableau[i]) == t.hidden[i] {
				t.hidden[i]--
			}
		}
		t.key = 'z'
	}
	t.Draw(s)
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
	s.Clear()
	for i, r := range qwerty {
		s.SetContent(6*i+1, 1, r, nil, black.Reverse(t.key == r))
		emitStr(s, 6*i+2, 1, black, "     ")
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
		s.SetContent(6*i+1, 6, r, nil, black.Reverse(t.key == r))
		emitStr(s, 6*i+2, 6, black, "     ")
	}
	for i := 0; i <= 6; i++ {
		for j := 0; j < t.hidden[i]; j++ {
			t.DrawCard(s, 6*i, 8+j, 52)
		}
		for j := Max(t.hidden[i], 0); j < len(t.tableau[i]); j++ {
			t.DrawCard(s, 6*i, 8+j, t.tableau[i][j])
		}
	}
	s.Sync()
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
			if t.IsWon() {
				close(quit)
				return
			}
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyCtrlC:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				case tcell.KeyRune:
					t.SelectKey(s, ev.Rune())
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	<-quit

	s.Fini()
}
