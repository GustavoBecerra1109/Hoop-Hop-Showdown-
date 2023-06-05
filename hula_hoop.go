package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Player represents a game player.
type Player struct {
	ID         int
	Team       int
	TokenCount int
}

// Game represents the game state and rules.
type Game struct {
	Hoops       []chan *Player
	Cones       []*sync.Mutex
	Tokens      []int
	NumTeams    int
	NumPlayers  int
	TotalTokens int
}

// NewGame initializes a new game.
func NewGame(numTeams, numPlayers int) *Game {
	game := &Game{
		Hoops:       make([]chan *Player, numTeams),
		Cones:       make([]*sync.Mutex, numTeams),
		Tokens:      make([]int, numTeams),
		NumTeams:    numTeams,
		NumPlayers:  numPlayers,
		TotalTokens: numTeams * 20, // Assuming 20 tokens per bucket
	}

	for i := 0; i < numTeams; i++ {
		game.Hoops[i] = make(chan *Player, 1)
		game.Cones[i] = &sync.Mutex{}
	}

	return game
}

// Run starts the game.
func (g *Game) Run() {
	var wg sync.WaitGroup
	wg.Add(g.NumPlayers)

	for i := 0; i < g.NumPlayers; i++ {
		go func(playerID int) {
			defer wg.Done()
			team := playerID % g.NumTeams
			player := &Player{ID: playerID, Team: team}

			for {
				hoop := <-g.Hoops[team] // Wait for available hoop

				if g.isWinner(player) {
					g.Cones[team].Lock()
					if g.Tokens[team] > 0 {
						g.Tokens[team]--
						player.TokenCount++
					}
					g.Cones[team].Unlock()
				}

				fmt.Printf("Player %d of Team %d: Jumping into hoop\n", player.ID, player.Team)
				time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Simulate jumping time

				if !g.isWinner(player) {
					fmt.Printf("Player %d of Team %d: Lost in RPS, stepping out of hoop\n", player.ID, player.Team)
					g.Hoops[team] <- hoop // Return hoop to available hoops
					return
				}

				if player.TokenCount == g.NumTeams {
					fmt.Printf("Player %d of Team %d: Won the game!\n", player.ID, player.Team)
					return
				}

				fmt.Printf("Player %d of Team %d: Reached the cone, getting a token\n", player.ID, player.Team)
				time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Simulate token retrieval time

				g.Cones[team].Lock()
				g.Tokens[team]++
				g.Cones[team].Unlock()

				g.Hoops[team] <- hoop // Return hoop to available hoops
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Game over.")
}

// isWinner simulates Rock, Paper, Scissors and determines if the player is a winner.
func (g *Game) isWinner(player *Player) bool {
	opponent := rand.Intn(g.NumPlayers)
	return player.ID != opponent
}

func main() {
	numTeams := 4
	numPlayers := numTeams

	game := NewGame(numTeams, numPlayers)

	for i := 0; i < numTeams; i++ {
		game.Hoops[i] <- &Player{} // Signal availability of hoop for each team
	}

	game.Run()
}
