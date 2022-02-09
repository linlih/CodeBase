package main

import (
	"fmt"
)

func main() {
	ak47, _ := getGun("ak47")
	musket, _ := getGun("musket")
	printDetail(ak47)
	printDetail(musket)
}

func printDetail(g iGun) {
	fmt.Printf("Gun: %s", g.getName())
	fmt.Printf(" Power: %d", g.getPower())
	fmt.Println()
}
