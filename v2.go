package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Player representa a un jugador del juego.
type Player struct {
	ID         int
	Team       int
	TokenCount int
	Positions  []string
}

// Game representa el estado y las reglas del juego.
type Game struct {
	Hoops       [][]chan *Player
	Tokens      []int
	NumTeams    int
	NumPlayers  int
	TotalTokens int
	Lock        sync.Mutex
	Winner      bool
}

// NewGame inicializa un nuevo juego.
func NewGame(numTeams, numPlayers int) *Game {
	game := &Game{
		Hoops:      make([][]chan *Player, numTeams),
		Tokens:     make([]int, numTeams),
		NumTeams:   numTeams,
		NumPlayers: numPlayers,
	}

	for i := 0; i < numTeams; i++ {
		game.Hoops[i] = make([]chan *Player, 5)
		for j := 0; j < 5; j++ {
			game.Hoops[i][j] = make(chan *Player, 1)
		}
	}

	return game
}

// Run inicia el juego.
func (g *Game) Run() {
	var wg sync.WaitGroup
	wg.Add(g.NumPlayers)

	for i := 0; i < g.NumPlayers; i++ {
		go func(playerID int) {
			defer wg.Done()
			team := playerID % g.NumTeams
			player := &Player{ID: playerID, Team: team}

			for {
				if g.Winner {
					return
				}

				hoop := <-g.Hoops[team][player.TokenCount%5] // Espera a que haya un aro disponible para el equipo

				g.Lock.Lock()
				if g.Tokens[team] == 20 {
					fmt.Printf("¡El equipo %d es el ganador!\n", team)
					g.Winner = true
					g.Lock.Unlock()
					go func() {
						g.Hoops[player.Team][player.TokenCount%5] <- hoop // Devuelve el aro a los aros disponibles
					}()
					return
				}
				g.Lock.Unlock()

				if player.TokenCount > 0 {
					player.Positions = append(player.Positions, fmt.Sprintf("[%d][%d]", player.Team, player.TokenCount%5))
				}

				if player.TokenCount%5 == 0 {
					// El jugador llega a la primera posición de otro equipo
					opponentTeam := (player.Team + 1) % g.NumTeams
					opponentIndex := 0

					fmt.Printf("¡El jugador del equipo %d llegó a la posición [%d][%d] del equipo %d!\n", player.Team, opponentTeam, opponentIndex, opponentTeam)
					time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Simula el tiempo de juego

					g.Lock.Lock()
					g.Tokens[player.Team]++
					player.TokenCount++
					player.Positions = append(player.Positions, fmt.Sprintf("[%d][%d]", opponentTeam, opponentIndex))
					g.Lock.Unlock()

					if g.Tokens[player.Team] == 20 {
						fmt.Printf("¡El equipo %d es el ganador!\n", player.Team)
						g.Lock.Lock()
						g.Winner = true
						g.Lock.Unlock()
						go func() {
							g.Hoops[player.Team][player.TokenCount%5] <- hoop // Devuelve el aro a los aros disponibles
						}()
						return
					}

					fmt.Printf("El jugador del equipo %d ganó un token y se mueve a la posición [%d][%d]!\n", player.Team, opponentTeam, opponentIndex)
					go func() {
						g.Hoops[opponentTeam][opponentIndex] <- hoop // Devuelve el aro al equipo perdedor
					}()
				} else {
					opponentTeam := hoop.Team
					opponentIndex := player.TokenCount % 5

					fmt.Printf("El jugador del equipo %d se encuentra con el jugador del equipo %d en la posición [%d][%d]\n", player.Team, opponentTeam, opponentTeam, opponentIndex)
					time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Simula el tiempo de juego

					for {
						if g.playRockPaperScissors(player, hoop) {
							g.Lock.Lock()
							g.Tokens[player.Team]++
							player.TokenCount++
							player.Positions = append(player.Positions, fmt.Sprintf("[%d][%d]", opponentTeam, opponentIndex))
							g.Lock.Unlock()

							if g.Tokens[player.Team] == 20 {
								fmt.Printf("¡El equipo %d es el ganador!\n", player.Team)
								g.Lock.Lock()
								g.Winner = true
								g.Lock.Unlock()
								go func() {
									g.Hoops[player.Team][player.TokenCount%5] <- hoop // Devuelve el aro a los aros disponibles
								}()
								return
							}

							fmt.Printf("El jugador del equipo %d le ha ganado al jugador del equipo %d!\n", player.Team, opponentTeam)
							go func() {
								g.Hoops[opponentTeam][opponentIndex] <- hoop // Devuelve el aro al equipo perdedor
							}()
							break
						} else {
							fmt.Printf("El jugador del equipo %d ha empatado contra el jugador del equipo %d\n", player.Team, opponentTeam)
							time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Simula el tiempo de juego
						}
					}
				}

				fmt.Printf("El equipo %d salta a la posición [%d][%d]\n", player.Team, player.Team, player.TokenCount%5)
				time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Simula el tiempo de salto

				if player.TokenCount%5 == 0 {
					newTeam := rand.Intn(g.NumTeams)
					fmt.Printf("El equipo %d cambió al equipo %d\n", player.Team, newTeam)
					player.Team = newTeam
				}
				go func() {
					g.Hoops[team][player.TokenCount%5] <- hoop // Devuelve el aro a los aros disponibles
				}()
			}
		}(i)
	}

	wg.Wait()
}

// playRockPaperScissors simula Piedra, Papel, Tijeras y determina si el jugador es el ganador.
func (g *Game) playRockPaperScissors(player *Player, hoop *Player) bool {
	for {
		playerChoice := rand.Intn(3) // 0 = Piedra, 1 = Papel, 2 = Tijeras
		hoopChoice := rand.Intn(3)

		fmt.Printf("El jugador del equipo %d elige %s, el jugador del equipo %d elige %s\n", player.Team, getChoiceName(playerChoice), hoop.Team, getChoiceName(hoopChoice))

		switch {
		case playerChoice == hoopChoice:
			fmt.Println("Empate, ambos jugadores eligen lo mismo")
		case playerChoice == 0 && hoopChoice == 2:
			fmt.Printf("El jugador del equipo %d gana con Piedra contra Tijeras\n", player.Team)
			return true // Piedra vence a Tijeras, el jugador gana
		case playerChoice == 1 && hoopChoice == 0:
			fmt.Printf("El jugador del equipo %d gana con Papel contra Piedra\n", player.Team)
			return true // Papel vence a Piedra, el jugador gana
		case playerChoice == 2 && hoopChoice == 1:
			fmt.Printf("El jugador del equipo %d gana con Tijeras contra Papel\n", player.Team)
			return true // Tijeras vence a Papel, el jugador gana
		default:
			fmt.Printf("El jugador del equipo %d pierde contra el jugador del equipo %d\n", player.Team, hoop.Team)
			return false // El jugador no gana
		}
	}
}

// getChoiceName devuelve el nombre correspondiente a la elección de Piedra, Papel o Tijeras.
func getChoiceName(choice int) string {
	switch choice {
	case 0:
		return "Piedra"
	case 1:
		return "Papel"
	case 2:
		return "Tijeras"
	default:
		return ""
	}
}

func main() {
	numTeams := 4
	numPlayers := numTeams

	game := NewGame(numTeams, numPlayers)

	for i := 0; i < numTeams; i++ {
		for j := 0; j < 5; j++ {
			go func(team, index int) {
				game.Hoops[team][index] <- &Player{} // Indica que hay un aro disponible para cada equipo
			}(i, j)
		}
	}

	game.Run()
}
