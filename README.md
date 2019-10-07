# Design Doc

A minimal Golang implementation of [Solitaire](https://bicyclecards.com/how-to-play/solitaire/) for the terminal, written as part of my application to [KP Engineering Fellows](https://fellows.kleinerperkins.com/).

**To run, download to your Go workspace and `go run solitaire.go`. You will need to have [Go](https://golang.org/doc/install) installed.**

Cards are repersented as integers 0-51, inclusive, and modded by 13 and 4 to get their rank and suit, respectively. All of the piles are repersented as arrays / slices of integers, with the tableau and foundations represented as arrays of slices.

[tcell](https://github.com/gdamore/tcell) is used to handle keyboard input and draw to the terminal.

The game uses pairs of consecutive key presses to move cards from one pile to another:

- **q** = Stock
- **w** = Talon
- **e, r, t, y** = Foundations
- **a, s, d, f, g, h, j** = Tableau

So, to move a card from Stock to the Talon, press **qw**.

I chose Go because I had never implemented anything non-trivial in it before and I wanted to learn.
