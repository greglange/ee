package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/greglange/ee/pkg/ee"
)

func main() {
	dominoes := ee.TwoDominoSets()
	hands := ee.DealDominoes(dominoes)

	fmt.Println(total_count + 18)

	fmt.Println("# show hands")
	for player, hand := range hands {
		fmt.Println("--- Player:", player, "---")
		fmt.Println(ee.DominoesString(hand))
	}

	reader := bufio.NewReader(os.Stdin)

	playerBid := -1
	var winningBid, trumpBid int

	fmt.Println("# do bids")
	for player, hand := range hands {
		for true {
			fmt.Println("--- Bid Player:", player, "---")
			if playerBid != -1 {
				fmt.Printf("Current bid: %d Trump: %d Player: %d\n", winningBid, trumpBid, playerBid)
			}
			fmt.Println(ee.DominoesString(hand))
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "p" || line == "pass" {
				break
			}
			var bid, trump int
			_, err := fmt.Sscanf(line, "%d %d", &bid, &trump)
			if err != nil {
				fmt.Println("Invalid input")
				continue
			}
			if playerBid == -1 {
				if bid >= 60 && bid <= 88 && trump >= -1 && trump <= 6 {
					if ee.HasSuit(trump, hand) {
						playerBid = player
						winningBid = bid
						trumpBid = trump
						break
					}
				}
			} else {
				if bid > winningBid && bid <= 88 && trump >= -1 && trump <= 6 {
					if ee.HasSuit(trump, hand) {
						playerBid = player
						winningBid = bid
						trumpBid = trump
						break
					}
				}
			}
			fmt.Println("Invalid bid")
		}
	}
	teamBid := playerBid % 2
	fmt.Printf("Winning bid: %d Trump: %d Team: %d Player: %d\n", winningBid, trumpBid, teamBid, playerBid)

	fmt.Println("# play hand")
	leadPlayer := playerBid
	tricks := [2][][]*ee.Domino{}
	for round := 0; round < 9; round++ {
		fmt.Println("+ round:", round)

		played := make(map[string]int)
		dominoesPlayed := []*ee.Domino{}
		currentPlayer := -1
		for currentPlayer != leadPlayer {
			if currentPlayer == -1 {
				currentPlayer = leadPlayer
			}
			var dominoIndex int
			var hand []*ee.Domino
			for true {
				fmt.Println("Player:", currentPlayer)
				fmt.Println("Played:", ee.DominoesString(dominoesPlayed))
				hand = hands[currentPlayer]
				fmt.Println("Hand:", ee.DominoesString(hand))
				line, _ := reader.ReadString('\n')
				_, err := fmt.Sscanf(line, "%d", &dominoIndex)
				if err != nil {
					fmt.Println("Invalid input")
					continue
				}
				if dominoIndex < 0 || dominoIndex >= len(hand) {
					fmt.Println("Invalid input")
					continue
				}
				if ee.ValidPlay(dominoIndex, hand, dominoesPlayed, trumpBid) {
					break
				}
			}
			var domino *ee.Domino
			domino, hands[currentPlayer] = ee.PlayDomino(hand, dominoIndex)
			played[domino.Id] = currentPlayer
			dominoesPlayed = append(dominoesPlayed, domino)
			currentPlayer = (currentPlayer + 1) % 6
		}
		winningDomino := ee.TrickWinner(dominoesPlayed, trumpBid)
		winningPlayer := played[winningDomino.Id]
		fmt.Println("Winning player:", winningPlayer)
		fmt.Println("Winning domino:", winningDomino)
		fmt.Println("Dominoes played:", ee.DominoesString(dominoesPlayed))
		tricks[winningPlayer%2] = append(tricks[winningPlayer%2], dominoesPlayed)
		leadPlayer = winningPlayer
	}

	fmt.Println("# tally scores")
	for team := 0; team < 2; team++ {
		fmt.Println("Tricks won by team:", team)
		for _, trick := range tricks[team] {
			fmt.Println(ee.DominoesString(trick))
		}
	}

	for team := 0; team < 2; team++ {
		fmt.Println("Team:", team)
		points := ee.CalcPoints(tricks[team])
		if teamBid == team {
			if winningBid <= points {
				fmt.Printf("   Team %d made their bid\n", team)
			} else {
				fmt.Printf("   Team %d did not make their bid\n", team)
			}
		}
		fmt.Printf("Team %d had %d points\n", team, points)
	}
}
