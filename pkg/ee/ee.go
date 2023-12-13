package ee

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const NORANK int = -1
const NOSUIT int = -1
const NOTRUMP int = -1

type Domino struct {
	Id     string
	EndOne int
	EndTwo int
}

func NewDomino(set int, endOne int, endTwo int) *Domino {
	one, two := endOne, endTwo
	if one < two {
		one, two = two, one
	}
	id := strings.Join([]string{strconv.Itoa(set), strconv.Itoa(one), strconv.Itoa(two)}, ".")
	return &Domino{id, one, two}
}

func (d *Domino) CountValue() int {
	if d.IsCount() {
		return d.Dots()
	}
	return 0
}

func (d *Domino) Dots() int {
	return d.EndOne + d.EndTwo
}

func (d *Domino) IsCount() bool {
	return d.Dots() > 0 && d.Dots()%5 == 0
}

func (d *Domino) IsDouble() bool {
	return d.EndOne == d.EndTwo
}

func (d *Domino) IsDoubleInSuit(suit int) bool {
	return d.EndOne == suit && d.EndOne == d.EndTwo
}

func (d *Domino) IsGreater(d2 *Domino, leadSuit int, trump int) bool {
	if d2.IsTrump(trump) {
		if !d.IsTrump(trump) {
			return false
		} else if d2.IsDouble() {
			return false
		} else if d.IsDouble() {
			return true
		} else {
			return d.Rank(trump) > d2.Rank(trump)
		}
	} else {
		if d.IsTrump(trump) {
			return true
		} else {
			if d2.IsDoubleInSuit(leadSuit) {
				return false
			} else if d.IsDoubleInSuit(leadSuit) {
				return true
			}
			return d.Rank(leadSuit) > d2.Rank(leadSuit)
		}
	}
}

func (d *Domino) IsSuit(suit int, trump int) bool {
	return !d.IsTrump(trump) && (d.EndOne == suit || d.EndTwo == suit)
}

func (d *Domino) IsTrump(trump int) bool {
	return d.EndOne == trump || d.EndTwo == trump
}

func (d *Domino) LeadSuit(trump int) int {
	if d.IsTrump(trump) {
		return trump
	}
	return d.EndOne
}

func (d *Domino) Rank(suit int) int {
	if d.EndOne == suit {
		return d.EndTwo
	} else if d.EndTwo == suit {
		return d.EndOne
	}
	return NORANK
}

func (d *Domino) String() string {
	return fmt.Sprintf("{%d:%d}", d.EndOne, d.EndTwo)
}

func DominoesString(dominoes []*Domino) string {
	dominoesString := ""
	for _, domino := range dominoes {
		if len(dominoesString) > 0 {
			dominoesString += " "
		}
		dominoesString += domino.String()
	}
	return dominoesString
}

func PlayDomino(dominoes []*Domino, i int) (*Domino, []*Domino) {
	domino := dominoes[i]
	return domino, append(dominoes[:i], dominoes[i+1:]...)
}

func DominoSet(set int) []*Domino {
	dominoes := []*Domino{}
	for i := 0; i < 7; i++ {
		for j := i; j < 7; j++ {
			dominoes = append(dominoes, NewDomino(set, i, j))
		}
	}
	return dominoes
}

func TwoDominoSets() []*Domino {
	dominoes := []*Domino{}
	for i := 0; i < 2; i++ {
		dominoes = append(dominoes, DominoSet(i)...)
	}
	return dominoes
}

func ShuffleDominoes(dominoes []*Domino) []*Domino {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(dominoes), func(i, j int) { dominoes[i], dominoes[j] = dominoes[j], dominoes[i] })
	return dominoes
}

func DealDominoes(dominoes []*Domino) [6][]*Domino {
	dominoes = ShuffleDominoes(dominoes)
	return [6][]*Domino{
		dominoes[0:9],
		dominoes[9:18],
		dominoes[18:27],
		dominoes[27:36],
		dominoes[36:45],
		dominoes[45:54],
	}
}

func CalcPoints(tricks [][]*Domino) int {
	points := len(tricks) * 2
	for _, trick := range tricks {
		for _, d := range trick {
			points += d.CountValue()
		}
	}
	return points
}

func TrickWinner(dominoes []*Domino, trump int) *Domino {
	winner := dominoes[0]
	leadSuit := winner.LeadSuit(trump)
	for _, domino := range dominoes[1:] {
		if domino.IsGreater(winner, leadSuit, trump) {
			winner = domino
		}
	}
	return winner
}

func HasSuit(suit int, hand []*Domino) bool {
	if suit == NOSUIT {
		return true
	}
	for _, d := range hand {
		if d.IsSuit(suit, NOTRUMP) {
			return true
		}
	}
	return false
}

func ValidPlay(index int, hand []*Domino, dominoesPlayed []*Domino, trump int) bool {
	if len(dominoesPlayed) == 0 {
		return true
	}

	lead := dominoesPlayed[0]
	domino := hand[index]

	if lead.IsTrump(trump) {
		if domino.IsTrump(trump) {
			return true
		}
		for _, d := range hand {
			if d.IsTrump(trump) {
				return false
			}
		}
		return true
	}

	suit := lead.EndOne
	if domino.IsSuit(suit, trump) {
		return true
	}
	for _, d := range hand {
		if d.IsSuit(suit, trump) {
			return false
		}
	}
	return true
}
