package tcpserver

import "fmt"

type Player struct {
	Id        int
	Hand      []int
	HoardView []int
	Money     int
}

func (player Player) ShowHoard() {
	fmt.Println(player.HoardView)
}

func (player Player) ShowHand() {
	fmt.Println(player.Hand)
}
