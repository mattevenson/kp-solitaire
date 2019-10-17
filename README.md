# Solitaire

_A minimal Golang implementation of [Solitaire](https://bicyclecards.com/how-to-play/solitaire/) for the terminal, written as part of my application to [KP Engineering Fellows](https://fellows.kleinerperkins.com/)._

**To run, download to your [Go](https://golang.org/doc/install) workspace and `go get github.com/gdamore/tcell && go run solitaire.go`.**

Cards are repersented as integers 0-51, inclusive, and modded by 13 and 4 to get their rank and suit, respectively.

All of the piles are repersented as arrays / slices of integers, with the tableau and foundations represented as arrays of slices. A seperate array stores the "frontier" of the tableau piles, with cards before that
frontier being hidden as denoted by `?`.

The game uses pairs of consecutive key presses to move cards from one pile to another:

- **q** = Stock
- **w** = Talon
- **e, r, t, y** = Foundations
- **a, s, d, f, g, h, j** = Tableau

_Ex: To move a card from the Stock to the Talon, press **qw**._

The core game logic just involves relating these letters to their corresponding piles,
and ensuring that a move is valid before it is executed.

[tcell](https://github.com/gdamore/tcell) is used to handle keyboard input and draw to the terminal.

I chose Go because I had never implemented anything non-trivial in it before and I wanted to learn.

[![asciicast](https://asciinema.org/a/n4xDFCB2WzJhUlmKcQKnZXGYc.svg)](https://asciinema.org/a/n4xDFCB2WzJhUlmKcQKnZXGYc)
