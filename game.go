package tcpserver

import (
	"fmt"
	"math/rand"
	"time"
)

type Game struct {
	Hoard        []int
	RandomSource *rand.Rand
	AntePrice    int
	AntePayer    Player
}

func CreateGame() Game {
	game := Game{}
	game.NewRandomSource()
	game.Hoard = make([]int, 0, 5)
	game.AntePrice = 10
	return game
}

func (game *Game) NewRandomSource() {
	game.RandomSource = rand.New(rand.NewSource(time.Now().UnixMicro()))
}

func (game *Game) RollD12() int {
	return game.RandomSource.Intn(12) + 1
}

func (game *Game) RollHoard() {
	for i := 0; i < 3; i++ {
		game.Hoard = append(game.Hoard, game.RollD12())
	}
}

func (game *Game) ShowHoard() {
	fmt.Println(game.Hoard)
}

func (game *Game) RollPlayerHand() []int {
	hand := make([]int, 2)
	hand[0] = game.RollD12()
	hand[1] = game.RollD12()
	return hand
}
